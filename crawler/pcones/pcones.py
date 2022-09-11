import math
import re
import time
import datetime
import json
from xml.dom.minidom import Attr

import requests
from bs4 import BeautifulSoup

from mq import MQ
from crawler import utils
import log_settings


logger = log_settings.get_logger(__name__)


def main():
    pcones = PconesCrawler()
    pcones.products_list_page_loop()


class PconesCrawler(object):

    def __init__(self, queue_name: str='mws'):
        self.url = 'https://www.1-s.jp/products/list/'
        self.query = dict(sorted({
            'mode': 'search',
            'size': 100,
            'pageno': 1,
            'name_op': 'AND',
        }.items(), key=lambda x: x[0]))
        self.max_pages = self.get_products_count()
        self.timestamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')
        self.mq = MQ(queue_name)

    def products_list_page_loop(self, interval_sec: int=2):
        logger.info('action=products_list_page_loop status=run')

        while self.query['pageno'] <= self.max_pages:
            logger.info(self.query['pageno'])
            response = utils.request(self.url, params=self.query)
            time.sleep(interval_sec)
            products = PconesHTMLPage.scrape_product_list_page(response.text)
            [self.publish_queue(product['jan'], product['cost']) for product in products]
            self.query['pageno'] += 1

        logger.info('action=products_list_page_loop status=done')

    def get_products_count(self, interval_sec: int=1) -> int:
        logger.info('action=get_products_count status=run')

        response = utils.request(self.url, params=self.query)
        time.sleep(interval_sec)
        products_count = PconesHTMLPage.scrape_products_count(response.text)
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
    def scrape_products_count(response: str)-> str:
        logger.info('action=scrape_products_count status=run')

        soup = BeautifulSoup(response, 'lxml')
        pagenumber = soup.select_one('.pagenumber').text

        logger.info('action=scrape_products_count status=done')
        return float(pagenumber)

    @staticmethod
    def scrape_product_list_page(response: str) -> list[dict]:
        logger.info('action=scrape_product_list_page status=run')
        
        products = []

        soup = BeautifulSoup(response, 'lxml')
        product_list = soup.select('#product_list tr')
        for product in product_list:
            list_code = product.select_one('.list_code')
            list_price = product.select_one('.list_price')

            if not list_code or not list_price:
                continue
            
            jan = re.search('[0-9]{13}', list_code.text)
            price = re.search('.*円', list_price.text)
            try:
                is_sold_out = product.select_one('.list_stock font').text
            except AttributeError as ex:
                is_sold_out = None

            if not jan or not price or is_sold_out == '×品切れ':
                continue
            else:
                price = int(''.join(re.findall('[0-9]', price.group())))
                products.append({'jan': jan.group(), 'cost': price})

        logger.info('action=scrape_product_list_page status=done')
        return products

    @staticmethod
    def scrape_product_detail_page(response: str) -> dict:

        PRODUCT_CODE_ROW_NUM = 2
        soup = BeautifulSoup(response, 'lxml')
        try:
            jan = soup.select('.detail_etc_table td')[PRODUCT_CODE_ROW_NUM].attrs.get('content')
            jan = re.fullmatch('[0-9]{13}', jan).group()
        except AttributeError as ex:
            logger.info({'message': 'jan code is None', 'error': ex})
            jan = None

        try:
            price = soup.select_one('.detail_price .price span').attrs.get('content')
            price = int(''.join(re.findall('[0-9]', price)))
        except AttributeError as ex:
            logger.info({'message': 'price is None', 'error': ex})
            price = None

        is_stocked = soup.select_one('#cart')

        return {'jan': jan, 'price': price, 'is_stocked': bool(is_stocked)}
