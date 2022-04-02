import math
import time
import datetime
import json

import requests
from bs4 import BeautifulSoup

from mq import MQ
from crawler import utils
import log_settings


logger = log_settings.get_logger(__name__)


class PconesCrawler(object):

    def __init__(self, queue_name: str='mws'):
        self.url = 'https://www.1-s.jp/products/list'
        self.query = {
            'mode': 'search',
            'orderby': 'date+desc',
            'size': 100,
            'pageno': 1,
        }
        self.max_pages = self.get_products_count()
        self.timestamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')
        self.mq = MQ(queue_name)

    def products_list_page_loop(self, interval_sec: int=2):
        logger.info('action=products_list_page_loop status=run')

        while self.query['pageno'] <= self.max_pages:
            response = utils.request(self.url, params=self.query)
            time.sleep(interval_sec)
            products = PconesHTMLPage.scrape_product_list_page(response)
            [self.publish_queue(product['jan'], product['cost']) for product in products]
            self.query['pageno'] += 1

        logger.info('action=products_list_page_loop status=done')

    def get_products_count(self, interval_sec: int=1) -> int:
        logger.info('action=get_products_count status=run')

        response = utils.request(self.url, params=self.query)
        time.sleep(interval_sec)
        products_count = PconesHTMLPage.scrape_products_count(response)
        max_pages = math.ceil(products_count / 100)

        logger.info('action=get_products_count')
        return max_pages

    def publish_queue(self, jan: str, cost: int):
        logger.info('action=publish_queue status=run')

        self.mq.publish(json.dumps({
            'filename': f'pcones_{self.timestamp}',
            'jan': jan,
            'cost': cost,
        }))
        logger.info('action=publish_queue status=done')


class PconesHTMLPage(object):

    @staticmethod
    def scrape_products_count(response: requests.Response) -> str:
        logger.info('action=scrape_products_count status=run')

        soup = BeautifulSoup(response.text, 'lxml')
        pagenumber = soup.select_one('.pagenumber').text

        logger.info('action=scrape_products_count status=done')
        return float(pagenumber)

    @staticmethod
    def scrape_product_list_page(response: requests.Response) -> list[dict]:
        logger.info('action=scrape_product_list_page status=run')
        pass







