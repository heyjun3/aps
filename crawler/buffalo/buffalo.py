import re
import logging
import time
import datetime
import os

import requests
import openpyxl
from requests import Response
from bs4 import BeautifulSoup

from crawler import utils
from crawler.models import BuffaloProduct
import settings

logger = logging.getLogger(__name__)


class Crawler():

    filename = 'buffalo'

    def __init__(self, url):
        self.db_list = []
        self.not_db_list = []
        self.url = url
        self.products = []
        self.session = requests.Session()

    def list_page_crawling(self):
        logger.info('action=list_page_crawling status=run')

        while True:
            response = utils.request(self.url, self.session)
            self.list_page_scraping(response)
            time.sleep(2)
            if not self.url:
                logger.info('action=list_page_crawling status=done')
                break

    def detail_page_crawling(self, products):
        logger.info('action=detail_page_crawling status=run')

        for product in products:
            response = utils.request(self.session, product.url)
            time.sleep(2)
            product.jan = self.detail_page_scraping(response)
            product.save()

        logger.info('action=detail_page_crawling status=done')

    def export_excel_file(self, products, save_path=settings.SCRAPE_SAVE_PATH):
        logger.info('action=export_excel_file status=run')
        dt = datetime.datetime.now()
        timestamp = dt.strftime('%Y%m%d_%H%M%S')

        workbook = openpyxl.Workbook()
        sheet = workbook['Sheet']
        sheet.append(['JAN', 'Cost'])

        for product in products:
            if product.jan and product.price:
                sheet.append([product.jan, product.price])
        
        workbook.save(os.path.join(save_path, f'{self.filename}{timestamp}.xlsx'))
        workbook.close()
        logger.info('action=export_excel_file status=done')

    def list_page_scraping(self, response: Response):
        logger.info('action=list_page_scraping status=run')

        soup = BeautifulSoup(response.text, 'lxml')
        products_list = soup.select_one('.list')
        products = products_list.select('li.clearfix')

        for product in products:
            title = product.select_one('h3').text
            title_flag = re.search('整備済|検査済', title)
            is_sold_out = product.select_one('.soldout')
            if title_flag or is_sold_out:
                continue
            price = product.select_one('.price span').text.strip()
            price = ''.join(re.findall('[\\d+]', price))
            url = product.select_one('.columnRight a').attrs['href']
            product_code = re.search('[\\d]+', url)
            buffalo = BuffaloProduct.create(name=title, price=int(price), url=settings.BUFFALO_URL + url,
                                            product_code=product_code.group())
            self.products.append(buffalo)
            logger.debug(buffalo.value)

        self.url = None
        logger.info('action=list_page_scraping status=done')

    @staticmethod
    def detail_page_scraping(response: Response):
        logger.info('action=detail_page_selector status=run')
        soup = BeautifulSoup(response.text, 'lxml')
        jan = soup.select_one('#detailBox02 .columnLeft p').get_text()
        jan = re.search('4[\\d]{12}', jan)
        if jan is None:
            return None
        return jan.group()

    def classify_exist_db(self):
        logger.info('action=classify_exist_db status=run')

        for product in self.products:
            db_response = BuffaloProduct.get_product_code(product.product_code)
            if not db_response:
                self.not_db_list.append(product)
            elif not db_response.price == product.price:
                BuffaloProduct.price_update('product_code', db_response.product_code, product.price)
                db_response.price = product.price
                self.db_list.append(db_response)
            else:
                self.db_list.append(db_response)

        logger.info('action=classify_exist_db status=done')


def main():
    start_url = settings.BUFFALO_START_URL
    crawler_buffalo = Crawler(url=start_url)
    crawler_buffalo.list_page_crawling()
    crawler_buffalo.classify_exist_db()
    crawler_buffalo.detail_page_crawling(crawler_buffalo.not_db_list)
    crawler_buffalo.export_excel_file(crawler_buffalo.db_list+crawler_buffalo.not_db_list,
                                      save_path=settings.SCRAPE_SCHEDULE_SAVE_PATH)
