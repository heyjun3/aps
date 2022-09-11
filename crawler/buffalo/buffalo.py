import re
import time
import datetime
import os
import json
from urllib.parse import urlparse
from urllib.parse import parse_qs

import requests
import openpyxl
from requests import Response
from bs4 import BeautifulSoup

from crawler import buffalo, utils
from crawler.buffalo.models import BuffaloProduct
import settings
import log_settings
from mq import MQ

logger = log_settings.get_logger(__name__)


class BuffaloCrawler():

    def __init__(self, start_url, queue_name: str = 'mws'):
        self.url = start_url
        self.buffalo_product_list = []
        self.mq = MQ(queue_name)
        self.timestamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')

    def pool_product_list_page(self, interval_sec: int = 2) -> None:
        logger.info('action=pool_product_list_page status=run')

        response = utils.request(url=self.url)
        self.buffalo_product_list.extend(BuffaloHTMLPage.scrape_product_list_page(response.text))
        time.sleep(interval_sec)

        logger.info('action=pool_product_list_page status=done')

    def pool_product_detail_page(self, interval_sec: int = 2) -> None:
        logger.info('action=pool_product_detail_page status=run')

        for buffalo_product in self.buffalo_product_list:
            product = BuffaloProduct.get_product_code(buffalo_product.product_code)
            if product is None:
                query = {'product_id': buffalo_product.product_code}
                response = utils.request(url=settings.BUFFALO_DETAIL_PAGE_URL, params=query)
                time.sleep(interval_sec)
                parsed_value = BuffaloHTMLPage.scrape_product_detail_page(response.text)
                buffalo_product.jan = parsed_value.get('jan')
                buffalo_product.save()
            else:
                buffalo_product.jan = product.jan

            if buffalo_product.jan:
                self.publish_queue(buffalo_product.jan, buffalo_product.price)

        logger.info('action=pool_product_detail_page status=done')

    def publish_queue(self, jan: str, price: int) -> None:
        logger.info('action=publish_queue status=run')

        self.mq.publish(json.dumps({
            'filename': f'buffalo_{self.timestamp}',
            'jan': jan,
            'cost': price,
        }))
        logger.info('action=publish_queue status=done')


class BuffaloHTMLPage(object): 

    @staticmethod
    def scrape_product_list_page(response: str) -> list[BuffaloProduct]:
        logger.info('action=scrape_product_list_page status=run')

        buffalo_product_list = []
        soup = BeautifulSoup(response, 'lxml')
        products_list = soup.select_one('.list')
        products = products_list.select('li.clearfix')

        for product in products:
            title = product.select_one('h3').text
            title_flag = re.search('整備済|検査済', title)
            is_sold_out = product.select_one('.soldout')
            if title_flag or is_sold_out:
                continue
            price = int(''.join(re.findall('[0-9]', product.select_one('.price span').text.strip())))
            url = product.select_one('.columnRight a').attrs['href']
            product_code = parse_qs(urlparse(url).query).get('product_id').pop()
            buffalo = BuffaloProduct(name=title, price=price, product_code=product_code)
            buffalo_product_list.append(buffalo)

        logger.info('action=scrape_product_list_page status=done')
        return buffalo_product_list

    @staticmethod
    def scrape_product_detail_page(response: str) -> dict|None:
        logger.info('action=scrape_product_detail_page status=run')
        soup = BeautifulSoup(response, 'lxml')
        jan = soup.select_one('#detailBox02 .columnLeft p').get_text()
        price = soup.select_one('#detailBox01 #price span').get_text()
        is_stocked = soup.select_one('#detailBox01 #cart')
        try:
            jan = re.search('[0-9]{13}', jan).group()
        except AttributeError as ex:
            logger.info(f'jan code is None error={ex}')
            jan = None
        
        price = ''.join(re.findall('[0-9]', price))
        
        return {'jan': jan, 'price': int(price), 'is_stocked': bool(is_stocked)}


def main():
    client = BuffaloCrawler(start_url=settings.BUFFALO_START_URL)
    client.pool_product_list_page()
    client.pool_product_detail_page()
