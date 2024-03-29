import re
import time
import math
import datetime
import json
from queue import Queue
from urllib.parse import urlparse
from urllib.parse import urljoin
from urllib.parse import parse_qsl
import threading
from typing import List
from typing import Callable
from copy import deepcopy
from collections import ChainMap
from functools import reduce
from functools import partial
from operator import itemgetter
from operator import attrgetter

import requests
from bs4 import BeautifulSoup
from first import first
from requests_html import HTMLSession

from crawler.rakuten.models import RakutenProduct
import settings
import log_settings
from mq import MQ
from crawler import utils


logger = log_settings.get_logger(__name__)


def logging(func):
    def _inner(*args, **kwargs):
        logger.info({'action': func.__name__, 'status': 'run'})
        result = func(*args, **kwargs)
        logger.info({'action': func.__name__, 'status': 'done'})
        return result
    return _inner


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
    rakuten_url = 'https://www.rakuten.co.jp/'
    PER_PAGE_COUNT = 45

    def __init__(self) -> None:
        self.timestamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')

    @logging
    def crawle_by_shop(self, shop_code: str):
        self.api_client = RakutenAPIClient(shop_code)
        self.shop_code = shop_code
        shop_id = self._get_shop_id(shop_code)
        query = {
            'sid': shop_id,
            'used': 0,
            's': 3,
            'p': 1,
            'max': 100000,
        }
        while query:
            query = self._search_sequence(query)

    @logging
    def _search_sequence(self, query: dict, interval_sec: int=1) -> dict|None:

        response = utils.request(url=self.base_url, params=query, time_sleep=interval_sec)
        logger.info({'request_url': response.url})
        parsed_value = RakutenHTMLPage.parse_product_list_page(response.text)
        parsed_value = [value | {'shop_code': self.shop_code} for value in parsed_value]
        rakuten_products = RakutenProduct.get_products_by_shop_code_and_product_codes(
            [value.get('product_code') for value in parsed_value], self.shop_code)
        products = self._mapping_rakuten_products(parsed_value, rakuten_products)

        search_products = [product for product in products if not 'jan' in product]
        responses = [utils.request(product.get('url'), time_sleep=interval_sec)
                                                    for product in search_products]
        searched_products = list(filter(None, reduce(lambda d, f: map(f, d), [
            RakutenHTMLPage.scrape_product_detail_page,
            partial(self._mapping_search_value, products=search_products),
        ], responses)))

        RakutenProduct.insert_all_on_conflict_do_nothing(searched_products)

        [self.api_client.mq.publish(value) for value in 
            reduce(lambda d, f: map(f, d), [
                self._calc_real_price,
                self._generate_enqueue_str,
            ], searched_products + [product for product in products if product.get('jan')]) 
            if value]

        if not parsed_value:
            return

        next_query = self._generate_next_page_query(response, parsed_value[-1])
        logger.info({'next page query': next_query})
        return next_query

    @logging
    def _get_shop_id(self, shop_code: str) -> int:
        res = utils.request(urljoin(self.rakuten_url, shop_code), time_sleep=1)
        res.html.render(timeout=60)
        shop_id = (shop_id if (shop_id := RakutenHTMLPage.parse_shop_id(res.html.html))
                    else shop_id if (shop_id := reduce(lambda d, f: f(d), [
                        RakutenHTMLPage.parse_header_html_src,
                        itemgetter('src'),
                        partial(urljoin, res.url),
                        partial(utils.request, time_sleep=1),
                        attrgetter('text'),
                        RakutenHTMLPage.parse_shop_id,
                        ], res.html.html)) else None)
        return shop_id

    @logging
    def _generate_next_page_query(self, response: requests.Response, last_product: dict) -> dict:
        next_query = reduce(lambda d, f: f(d), [
                    RakutenHTMLPage.parse_next_page_url,
                    itemgetter('url'),
                    urlparse,
                    attrgetter('query'),
                    parse_qsl,
                    dict
                    ],response)
        if next_query:
            return next_query
        
        current_query = reduce(lambda d, f: f(d), [
            urlparse,
            attrgetter('query'),
            parse_qsl,
            dict,
        ], response.url)

        if not current_query.get('p') == '150':
            return

        price = last_product.get('price')
        if not price:
            logger.error({'message': "product hasn't price Exception"})
            raise Exception

        return current_query | {'p': '1', 'max': int(float(price)-1)}

    @logging
    def _mapping_rakuten_products(self, parsed_products: List[dict], 
                                        rakuten_products: List[RakutenProduct]) -> List[dict]:
        
        products = [product | {'jan': rakuten_product.jan} 
            if (rakuten_product := first(rakuten_products, key=lambda x: x.product_code == product['product_code']))
            else (product | {'jan': jan.group()}) 
            if (jan := re.fullmatch('[0-9]{13}', product['product_code'])) 
            else product for product in parsed_products]

        return products

    @logging
    def _mapping_search_value(self, search_value: dict, products: List[dict]) -> dict:

        product = first(products, key=lambda x: x['product_code'] == search_value['product_code'])
        return product | {'jan': search_value.get('jan')} if product else None

    @logging
    def _calc_real_price(self, value: dict, discount_rate: float=0.9) -> dict|None:
        real_value = deepcopy(value)
        price = real_value.get('price')
        point = real_value.get('point')
        if not all((price, point)):
            logger.error({'message': 'Bad Parameter', 'value': value})
            return

        real_value['price'] = int((price * discount_rate) - point)

        return real_value

    @logging
    def _generate_enqueue_str(self, value: dict) -> str|None:

        if value is None:
            logger.error({'message': 'publish queue bad parameter', 'value': value})
            return 

        jan = value.get('jan')
        price = value.get('price')
        url = value.get('url')
        if not all((jan, price, url)):
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
    def scrape_product_detail_page(response: requests.Response|str) -> dict:
        logger.info('action=scrape_product_detail_page status=run')
        PRODUCT_CODE_INDEX = -1

        text = response if isinstance(response, str) else response.text
        soup = BeautifulSoup(text, 'lxml')
        jan = ((jan.group() if (jan := re.fullmatch('[0-9]{13}', code.get('value'))) 
                            else None) if (code := soup.select_one('#ratRanCode')) else None)
        if jan is None:
            jan = jan.get("content") if (jan := soup.select_one('meta[itemprop="gtin13"]')) else None
            if jan is not None:
                jan = match.group() if (match := re.fullmatch("[0-9]{13}", jan)) else None

        price = int(''.join(re.findall('[0-9]', price.text))) \
                                         if (price := soup.select_one('.price2')) else None

        if price is None:
            price = price_tag.get("content") if (price_tag := soup.select_one('meta[itemprop="price"]')) else None
            if price is not None:
                price = int("".join(re.findall("[0-9]", price)))

        url = (tag.get('content') 
              if (tag := soup.select_one('meta[property="og:url"]')) else 
              tag.get('href') 
              if (tag := soup.select_one('link[rel="canonical"]')) else 
              None)
        product_code = list(filter(None, urlparse(url).path.split('/')))[PRODUCT_CODE_INDEX] \
                                                                if url else None
        is_stocked = bool(soup.select_one('.cart-button-container'))
        if not is_stocked:
            app_data = soup.select_one("#item-page-app-data")
            if app_data:
                data = json.loads(app_data.text)
                try:
                    quantity = data["api"]["data"]["itemInfoSku"]["purchaseInfo"]["newPurchaseSku"]["quantity"]
                    is_stocked = bool(float(quantity))
                except (KeyError, ValueError)as e:
                    logger.error(e)
        
        logger.info('action=scrape_product_detail_page status=done')
        return {'jan': jan,
                'price': price,
                'product_code': product_code,
                'is_stocked': is_stocked}

    @staticmethod
    def parse_product_list_page(response: str) -> List[dict]:
        logger.info({'action': 'parse_product_list_page', 'status': 'run'})

        result = []

        soup = BeautifulSoup(response, 'lxml')
        products = soup.select('.searchresultitem')
        for product in products:
            name = product.select_one('.content.title a')
            url = product.select_one('.image a')
            price = product.select_one('.price--OX_YW')
            point = product.select_one('.points--AHzKn span')
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

    @staticmethod
    def parse_shop_id(response: str) -> int:
        soup = BeautifulSoup(response, 'lxml')

        sid = (sid.get('value') if (sid := soup.select_one('input[name=sid]')) else None)
        if sid is None:
            logger.error({"action": "parse_shop_id", 'message': 'sid is None'})
            return

        if not re.fullmatch('[0-9]+', str(sid)):
            logger.error({'message': 'sid validation is failed', 'value': sid})
            return

        return int(sid)

    @staticmethod
    def parse_next_page_url(response: requests.Response|str) -> dict:

        text = response if isinstance(response, str) else response.text
        soup = BeautifulSoup(text, 'lxml')
        url = tag.get('href') if (tag := soup.select_one('.nextPage')) else None
        return {'url': url}

    @staticmethod
    def parse_header_html_src(response: requests.Response|str) -> dict:
        text = response if isinstance(response, str) else response.text
        soup = BeautifulSoup(text, 'lxml')
        header_src = elem.get('src') if (elem := soup.select_one('.header-out')) else None
        if header_src is None:
            logger.error({"action": "parse_header_html_src",
                          "message": 'header html src is None'})
            raise Exception
        return {'src': header_src}


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
