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

    def pool_rakuten_request_api(self, interval_sec: int = 2):
        logger.info('action=main status=run')
        
        while True:
            logger.info(self.params)
            response = utils.request(url=settings.REQUEST_URL, params=self.params)
            time.sleep(interval_sec)

            rakuten_product_list = RakutenAPIJSON.get_rakuten_products(response.json())
            for rakuten_product in rakuten_product_list: 
                rakuten_product.shop_code = self.shop_code
                self.rakuten_product_queue.put(rakuten_product)

            if len(rakuten_product_list) < 30:
                self.rakuten_product_queue.put(None)
                break

            if self.params['page'] == 100:
                self.params['page'] = 1
                last_product_price = rakuten_product_list.pop().price
                if self.params['maxPrice'] == last_product_price:
                    last_product_price -= 100
                self.params['maxPrice'] = last_product_price
            else:
                self.params['page'] += 1

    def pool_rakuten_product_detail_page(self, interval_sec: int = 2):
        logger.info('action=pool_rakuten_product_detail_page status=run')

        while True:
            rakuten_product = self.rakuten_product_queue.get()
            if rakuten_product is None:
                break
            
            if not rakuten_product.jan:
                response = RakutenProduct.get_object_filter_productcode_and_shopcode(rakuten_product.product_code, rakuten_product.shop_code)
                if response is None:
                    url = urljoin(settings.RAKUTEN_ENDPOINT, f'{rakuten_product.shop_code}/{rakuten_product.product_code}')
                    response = utils.request(url=url)
                    time.sleep(interval_sec)
                    rakuten_product.jan = RakutenHTMLPage.scrape_product_detail_page(response.text)
                else:
                    rakuten_product.jan = response.jan

            self.publish_queue(rakuten_product.jan, rakuten_product.price)
            
        logger.info('action=pool_rakuten_product_detail_page status=done')

    def publish_queue(self, jan: str, price: int) -> None:
        logger.info('action=publish_queue status=run')

        if not jan or not price:
            return None

        params = {
            'filename': f'rakuten_{self.timestamp}',
            'jan': jan,
            'price': price,
        }
        self.mq.publish(json.dumps(params))
        
        logger.info('action=publish_queue status=done')
        return None


class RakutenHTMLPage(object):

    @staticmethod
    def scrape_product_detail_page(response: str) -> str|None:
        logger.info('action=scrape_product_detail_page status=run')

        soup = BeautifulSoup(response, 'lxml')
        try:
            jan = re.fullmatch('[\d]{13}', soup.select_one('#ratRanCode').get('value'))
        except AttributeError as e:
            logger.error(f'{e}')
            return None
        
        logger.info('action=scrape_product_detail_page status=done')
        return jan


class RakutenAPIJSON(object):

    @staticmethod
    def get_rakuten_products(response: dict) -> list[RakutenProduct]:
        logger.info('action=get_rakuten_products status=run')

        rakuten_product_list = []

        for item in response['Items']:
            price = item['Item']['itemPrice']
            point_rate = item['Item']['pointRate']
            price = int(int(price) * (91 - int(point_rate)) / 100)
            item_name = item['Item']['itemName']
            jan = RakutenAPIJSON.get_jan_code(item)
            rakuten_product = RakutenProduct(name=item_name, jan=jan, price=price)
            rakuten_product_list.append(rakuten_product)

        logger.info('action=get_rakuten_products status=done')
        return rakuten_product_list

    @staticmethod
    def get_jan_code(item: dict) -> str|None:
        logger.info('action=get_jan_code status=run')

        item_url = item['Item']['itemUrl']
        jan = re.search('[0-9]{13}', item_url)

        if jan is None:
            item_caption = item['Item']['itemCaption']
            jan = re.search('[0-9]{13}', item_caption)

        logger.info('action=get_jan_code status=done')
        return jan
