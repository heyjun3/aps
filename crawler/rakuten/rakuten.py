import re
import time
import math
import datetime
import json
from queue import Queue
from urllib.parse import urljoin
import threading
from typing import List
from typing import Callable
from copy import deepcopy
from collections import ChainMap
from functools import reduce

import requests
from bs4 import BeautifulSoup

from crawler.rakuten.models import RakutenProduct
import settings
import log_settings
from mq import MQ
from crawler import utils


logger = log_settings.get_logger(__name__)


class RakutenAPIClient:
    """Rakuten api Class"""
    def __init__(self, shop_code: str, queue_name: str = 'mws'):
        self.shop_code = shop_code
        self.rakuten_product_queue = Queue()
        self.mq = MQ(queue_name)
        self.timestamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')
        self.params = {
            'applicationId': settings.RAKUTEN_APP_ID,
            'shopCode': self.shop_code,
            'page': 1,
            'sort': '-itemPrice',
            'maxPrice': None,
        }

    def run_rakuten_search(self):
        logger.info('action=run_rakuten_search status=run')

        thread = threading.Thread(target=self.pool_rakuten_product_detail_page)
        thread.start()
        self.pool_rakuten_request_api()
        thread.join()

        logger.info('action=run_rakuten_search status=done')

    def pool_rakuten_request_api(self, interval_sec: int=2) -> None:
        logger.info('action=pool_rakuten_request_api status=run')
        
        while True:
            logger.info(self.params)
            response = utils.request(url=settings.REQUEST_URL, params=self.params)
            time.sleep(interval_sec)

            rakuten_products = RakutenAPIJSON.parse_rakuten_api_json(response.json())
            for rakuten_product in rakuten_products: 
                rakuten_product['shop_code'] = self.shop_code
                self.rakuten_product_queue.put(rakuten_product)

            if len(rakuten_products) < 30:
                self.rakuten_product_queue.put(None)
                break

            if self.params['page'] == 100:
                self.params['page'] = 1
                last_product_price = rakuten_products[-1]['price']
                if self.params['maxPrice'] == last_product_price:
                    last_product_price -= 100
                self.params['maxPrice'] = last_product_price
                continue

            self.params['page'] += 1

    def pool_rakuten_product_detail_page(self, interval_sec: int = 2):
        logger.info('action=pool_rakuten_product_detail_page status=run')

        while True:
            product = self.rakuten_product_queue.get()
            if product is None:
                break
            
            if not product.get('jan'):
                rakuten_product = RakutenProduct.get_object_filter_productcode_and_shopcode(product['product_code'], product['shop_code'])
                if rakuten_product is None:
                    response = utils.request(url=product.get('url'))
                    time.sleep(interval_sec)
                    parsed_value = RakutenHTMLPage.scrape_product_detail_page(response.text)
                    product['jan'] = parsed_value.get('jan')
                    RakutenProduct(
                        name=product['name'],
                        jan=product['jan'],
                        price=product['price'],
                        shop_code=product['shop_code'],
                        product_code=product['product_code'],
                        url=product['url'],
                    ).save()
                else:
                    product['jan'] = rakuten_product.jan

            self.publish_queue(product['jan'], product['price'], product['url'])
            
        logger.info('action=pool_rakuten_product_detail_page status=done')

    def publish_queue(self, jan: str, price: int, url: str) -> None:
        logger.info('action=publish_queue status=run')

        if not all([jan, price, url]):
            return

        self.mq.publish(json.dumps({
            'filename': f'rakuten_{self.timestamp}',
            'jan': jan,
            'cost': price,
            'url': url,
        }))
        
        logger.info('action=publish_queue status=done')


class RakutenCrawler(object):

    base_url = 'https://search.rakuten.co.jp/search/mall/'
    PER_PAGE_COUNT = 45

    def __init__(self, shop_id: str, shop_code: str) -> None:
        self.shop_id = shop_id
        self.shop_code = shop_code
        self.timestamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')
        self.api_client = RakutenAPIClient(shop_code)
        self.query = {
            'sid': shop_id,
            'used': 0,
            's': 3,
            'p': 1,
        }

    def main(self):
        
        response = utils.request(self.base_url, params=self.query)
        max_count = self._get_max_page_count(response)
        querys = self._create_querys(max_count)
        for query in querys:
            self.search_sequence(query)

    def search_sequence(self, query: dict, interval_sec: int=1):
        logger.info({'action': 'search_sequence', 'status': 'run'})

        response = utils.request(url=self.base_url, params=query)
        time.sleep(interval_sec)
        parsed_value = RakutenHTMLPage.parse_product_list_page(response.text)
        publish_product, search_products = self._mapping_rakuten_products(parsed_value)
        searched_products = reduce(lambda d, f: map(f, d), [
            self._search_detail_page,
            self._validate_searched_products], search_products)
        searched_products = list(filter(None, search_products))
        RakutenProduct.insert_all_on_conflict_do_nothing(searched_products)

        list(reduce(lambda d, f: map(f, d), [
            self._calc_real_price,
            self._generate_enqueue_str,
            self.api_client.mq.publish,
        ], publish_product + searched_products))

        logger.info({'action': 'search_sequence', 'status': 'done'})

    def _get_max_page_count(self, response: requests.Response) -> int:
        logger.info({'action': 'get_max_page_count', 'status': 'run'})

        max_products_count = RakutenHTMLPage.parse_max_products_count(response.text)
        max_page_count = math.ceil(max_products_count / self.PER_PAGE_COUNT)

        logger.info({'action': 'get_max_page_count', 'status': 'done'})
        return max_page_count

    def _create_querys(self, max_count: int) -> List[dict]:
        logger.info({'action': 'create_querys', 'status': 'run'})

        querys = [self.query | {'p': i+1} for i in range(max_count)]
            
        logger.info({'action': 'create_querys', 'status': 'done'})
        return querys

    def _mapping_rakuten_products(self, values: List[dict]) -> tuple[List[str], List[str]]:
        logger.info({'action': 'mapping_rakuten_products', 'status': 'run'})

        product_codes = [value.get('product_code') for value in values]
        rakuten_products = RakutenProduct.get_products_by_shop_code_and_product_codes(product_codes, self.shop_code)
        rakuten_products = [{product.product_code: product} for product in rakuten_products]
        rakuten_products = ChainMap(*rakuten_products)

        search_products = []
        for product_code in product_codes:
            rakuten_product = rakuten_products.get(product_code['product_code'])
            product_code.setdefault('shop_code', self.shop_code)
            if rakuten_product:
                product_code['jan'] = rakuten_product.jan
            else:
                search_products.append(product_code)
        product_codes = list(filter(lambda x: x.get('jan') is not None, product_codes))

        return product_codes, search_products

    def _search_detail_page(self, value: dict, interval_sec: int=1) -> dict:
        logger.info({'action': 'search_detail_page', 'status': 'run'})

        result = deepcopy(value)

        response = utils.request(value.get('url'))
        time.sleep(interval_sec)
        parsed_value = RakutenHTMLPage.scrape_product_detail_page(response.text)
        result['jan'] = parsed_value.get('jan')

        logger.info({'action': 'search_detail_page', 'status': 'run'})
        return result

    def _validate_searched_products(self, value: dict) -> dict|None:
        required = ('name', 'jan', 'price', 'shop_code', 'product_code', 'url')

        if not all([key in value for key in required]):
            logger.error({'message': 'validation fail Badparameter', 'value': value})
            return

        return value

    def _calc_real_price(self, value: dict, discount_rate: float=0.9) -> dict|None:
        real_value = deepcopy(value)
        price = real_value.get('price')
        point = real_value.get('point')
        if not all((price, point)):
            logger.error({'message': 'Bad Parameter', 'value': value})
            return
        real_price = int((price * discount_rate) - point)
        real_value['price'] = real_price

        return real_value

    def _generate_enqueue_str(self, value: dict) -> str|None:

        jan = value.get('jan')
        price = value.get('price')
        url = value.get('url')
        if not all((jan, price, url)) or value is None:
            logger.error({'message': 'publish queue bad parameter', 'value': value})
            return

        return json.dumps({
            'cost': price,
            'filename': f'rakuten_{self.timestamp}',
            'jan': jan,
            'url': url,
        })


class RakutenHTMLPage(object):

    @staticmethod
    def scrape_product_detail_page(response: str) -> dict:
        logger.info('action=scrape_product_detail_page status=run')

        soup = BeautifulSoup(response, 'lxml')
        try:
            jan = re.fullmatch('[0-9]{13}', soup.select_one('#ratRanCode').get('value')).group()
        except AttributeError as e:
            logger.error(f'{e}')
            return {}

        try:
            price = int(''.join(re.findall('[0-9]', soup.select_one('.price2').text)))
        except (AttributeError, TypeError) as ex:
            logger.info(f'price is None message {ex}')
            price = None

        is_stocked = soup.select_one('.cart-button-container')
        
        logger.info('action=scrape_product_detail_page status=done')
        return {'jan': jan,
                'price': price,
                'is_stocked': bool(is_stocked)}

    @staticmethod
    def parse_product_list_page(response: str) -> dict:
        logger.info({'action': 'parse_product_list_page', 'status': 'run'})

        result = []

        soup = BeautifulSoup(response, 'lxml')
        products = soup.select('.searchresultitem')
        for product in products:
            name = product.select_one('.content.title a')
            url = product.select_one('.image a')
            price = product.select_one('.important')
            point = product.select_one('.content.points span')
            if not all((name, url, price, point)):
                logger.error({
                    "message": 'parse value not Found Error',
                    'action': 'parse_product_list_page',
                    'parameters': {'name': name, 'url': url, 'price': price, 'point': point}})
                continue

            url = url.attrs.get('href')
            price = int(''.join(re.findall('[0-9]', price.text)))

            try:
                point = int(re.match('[0-9]+', point.text.replace(',', '')).group())
            except (ValueError, TypeError) as ex:
                logger.error({
                    'messages': f'point is Bad value error={ex}',
                    'point': point})
                continue               

            try:
                product_code = url.split('/')[-2]
            except IndexError as ex:
                logger.error({'messages': ex, 'action': 'parse_product_list_page'})
                continue

            result.append({
                'name': name.text,
                'price': price,
                'product_code': product_code,
                'point': point,
                'url': url,
            })
        
        logger.info({'action': 'parse_product_list_page', 'status': 'done'})
        return result

    @staticmethod
    def parse_max_products_count(response: str) -> int:
        logger.info({'action': 'parse_max_products_count', 'status': 'run'})

        soup = BeautifulSoup(response, 'lxml')
        products_count = soup.select_one('._medium')
        if not products_count:
            logger.error({'action': 'parse_max_products_count', 'message': 'html page has not products count'})
            raise MaxProductsCountNotFoundException
        
        counts = re.findall('[0-9]+', products_count.text.replace(',', ''))
        max_count = max(map(int, counts))

        logger.info({'action': 'parse_max_products_count', 'status': 'done'})
        return max_count


class RakutenAPIJSON(object):

    @staticmethod
    def parse_rakuten_api_json(response: dict) -> list[dict]:
        logger.info('action=parse_rakuten_api_json status=run')

        rakuten_products = []

        for item in response['Items']:
            price = item['Item']['itemPrice']
            point_rate = item['Item']['pointRate']
            price = int(int(price) * (91 - int(point_rate)) / 100)
            item_name = item['Item']['itemName']
            jan = RakutenAPIJSON.get_jan_code(item)
            item_code = item['Item']['itemCode'].split(':')
            product_code = item_code.pop()
            shop_code = item_code.pop()
            url = item['Item']['itemUrl']
            rakuten_products.append({
                'name': item_name,
                'jan': jan,
                'price': price,
                'product_code': product_code,
                'shop_code': shop_code,
                'url': url,
            })

        logger.info('action=parse_rakuten_api_json status=done')
        return rakuten_products

    @staticmethod
    def get_jan_code(item: dict) -> str|None:
        logger.info('action=get_jan_code status=run')

        item_url = item['Item']['itemUrl']
        jan = re.search('[0-9]{13}', item_url)

        if jan is None:
            item_caption = item['Item']['itemCaption']
            jan = re.search('[0-9]{13}', item_caption)

        if jan:
            logger.info('action=get_jan_code status=done')
            return jan.group()
        return


class MaxProductsCountNotFoundException(Exception):
    pass
