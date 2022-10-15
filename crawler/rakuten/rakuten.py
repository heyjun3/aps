import re
import time
import datetime
import json
from queue import Queue
from urllib.parse import urljoin
import threading

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
        return {'jan': jan, 'price': price, 'is_stocked': bool(is_stocked)}

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
            if not all((name, url, price)):
                logger.error({
                    "message": 'parse value not Found Error',
                    'action': 'parse_product_list_page',
                    'parameters': {'name': name, 'url': url, 'price': price}})
                continue

            url = url.attrs.get('href')
            price = int(''.join(re.findall('[0-9]', price.text)))
            try:
                product_code = url.split('/')[-2]
            except IndexError as ex:
                logger.error({'messages': ex, 'action': 'parse_product_list_page'})
                continue

            result.append({
                'name': name.text,
                'url': url,
                'price': price,
                'product_code': product_code,
            })
        
        logger.info({'action': 'parse_product_list_page', 'status': 'done'})
        return result


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
