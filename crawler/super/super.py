import re
import time
import urllib
from datetime import datetime
import json

from requests import Session
import requests
from bs4 import BeautifulSoup
import pandas as pd

import settings
import log_settings
from mq import MQ
from crawler import utils
from crawler.super.models import SuperProduct
from crawler.super.models import SuperShop
from crawler.super.models import SuperProductDetails


logger = log_settings.get_logger(__name__)


class SuperCrawler(object):

    def __init__(self, url: str, params: dict = None, timestamp: datetime = datetime.now(), queue_name: str = 'mws'):
        self.url = requests.Request(method='GET', url=url, params=params).prepare().url
        self.session = self.login()
        self.super_product_list = []
        self.timestamp = timestamp.strftime('%Y%m%d_%H%M%S')
        self.mq = MQ(queue_name)
        self.favorite_product_list = []

    def login(self) -> Session:
        logger.info('action=login status=run')

        session = requests.session()
        response = utils.request(url=settings.SUPER_LOGIN_URL, method='POST', session=session, data=settings.SUPER_LOGIN_INFO)
        
        return session

    def start_search_products(self) -> None:
        logger.info('action=start_search_products status=run')

        self.pool_product_list_page()
        self.pool_product_detail_page()

        logger.info('action=start_search_products status=done')

    def pool_shop_list_page(self, interval_sec: int = 2) -> None:
        logger.info('action=pool_shop_list_page status=run')

        while self.url is not None:
            logger.info(self.url)
            response = utils.request(url=self.url, session=self.session)
            super_shop_list = SuperHTMLPage.scrape_shop_list_page(response.text)
            self.url = SuperHTMLPage.scrape_next_page_url(response.text)
            [super_shop.save() for super_shop in super_shop_list]
            time.sleep(interval_sec)

        logger.info('action=pool_shop_list_page status=done')

    def pool_product_list_page(self, interval_sec: int = 2) -> None:
        logger.info('action=pool_product_list_page status=run')

        while self.url is not None:
            logger.info(self.url)
            response = utils.request(url=self.url, session=self.session)
            time.sleep(interval_sec)
            self.super_product_list.extend(SuperHTMLPage.scrape_product_list_page(response.text, response.url))
            self.url = SuperHTMLPage.scrape_next_page_url(response.text)
        
        logger.info('action=pool_product_list_page status=done')


    def pool_product_detail_page(self, interval_sec: int = 2):
        logger.info('action=get_product_detail_page status=run')

        for super_product in self.super_product_list:
            products = SuperProductDetails.get(super_product.product_code, super_product.price)
            if not products:
                url = urllib.parse.urljoin(settings.SUPER_DOMAIN_URL, f'p/r/pd_p/{super_product.product_code}')
                response = utils.request(url=url, session=self.session)
                time.sleep(interval_sec)
                super_product_details_list = SuperHTMLPage.scrape_product_detail_page(response.text)
                [self.publish_queue(product.jan, product.price) for product in super_product_details_list]

        logger.info('action=get_product_detail_page status=done')

    def pool_favorite_product_list_page(self, interval_sec: int = 2):
        logger.info('action=pool_favorite_product_list_page status=run')


        while self.url is not None:
            response = utils.request(url=self.url, session=self.session)
            time.sleep(interval_sec)
            self.favorite_product_list.extend(SuperHTMLPage.scrape_favorite_product_list_page(response.text))
            self.url = SuperHTMLPage.scrape_next_page_url(response.text)
        
        df = pd.DataFrame(data=None, columns={'jan': str, 'cost': int})
        for item in self.favorite_product_list:
            url, cost = item
            product_code = url.split('/')[-2]
            products = SuperProductDetails.get_objects_to_product_code(product_code)

            if products is None:
                response = utils.request(url=url, session=self.session)
                time.sleep(interval_sec)
                products = SuperHTMLPage.scrape_product_detail_page(response.text)

            for product in products:
                    if product.jan:
                        df = df.append({'jan': product.jan, 'cost': cost}, ignore_index=True)
                    else:
                        logger.info(product.value)

        df = df.dropna()
        return df

    def publish_queue(self, jan: str, price: int) -> None:
        logger.info('action=publish_queue status=run')

        if not jan or not price:
            return None
        params = {
            'filename': f'super_{self.timestamp}',
            'jan': jan,
            'cost': price,
        }
        self.mq.publish(json.dumps(params))

        logger.info('action=publish_queue status=done')


class SuperHTMLPage(object):

    @staticmethod
    def scrape_product_list_page(response: str, response_url: str, sales_tax: float = 1.1) -> list[SuperProduct]:
        logger.info('action=scrape_product_list_page status=run')

        super_product_list = []

        soup = BeautifulSoup(response, 'lxml')
        products = soup.select('.itembox-parts')
        for product in products:
            try:
                item_name = product.select_one('.item-name a')
                name = item_name.text.strip().replace('\u3000', '')
                product_code = re.search('[\d]+', item_name.attrs.get('href')).group()
                shop_code = re.search('[\d]+', urllib.parse.urlparse(response_url).path).group()
                price = product.select_one('.item-price').text
                price = int(int(''.join(re.findall('\\d+', price))) * sales_tax)
                item = SuperProduct(name=name, product_code=product_code, shop_code=shop_code, price=price)
                item.save()

            except AttributeError as e:
                logger.error(f'action=scrape_product_list_page error={e}')
                continue
            super_product_list.append(item)

        logger.info('action=scrape_product_list_page status=done')
        return super_product_list

    @staticmethod
    def scrape_product_detail_page(response: str, consume_tax_rate: float = 1.1) -> list[SuperProductDetails]:
        logger.info('action=scrape_product_detail_page status=run')
        super_detail_product_list = []

        soup = BeautifulSoup(response, 'lxml')
        table = soup.select('.ts-tr02')
        product_code = re.search('[\d]+', soup.select_one('.co-fs12.co-clf.reduce-tax .co-pc-only').text).group()
        shop_code = re.search('[\d]+', urllib.parse.urlparse(soup.select_one('.dl-name-txt').get('href')).path).group()

        for row in table:
            try:
                jan = re.search('[0-9]{13}', row.select_one('.co-fcgray.td-jan').text).group()
            except AttributeError as e:
                logger.error(f"product hasn't jan code error={e}")
                jan = None
            price = int(int(''.join(re.findall('[\d]+', row.select_one('.td-price02').text))) * consume_tax_rate)
            set_number = int(re.search('[\d]+', row.select_one('.co-align-center.co-pc-only.border-rt.border-b').text.strip()).group())
            super_product = SuperProductDetails(jan=jan, price=price, set_number=set_number, shop_code=shop_code, product_code=product_code)
            super_detail_product_list.append(super_product)

        logger.info('action=scrape_product_detail_page status=done')
        return super_detail_product_list

    @staticmethod
    def scrape_shop_list_page(response: str) -> list[SuperShop]:
        logger.info('action=scrape_shop_list_page status=run')
        super_shop_list = []

        soup = BeautifulSoup(response, 'lxml')
        shop_list = soup.select('.dealer-eachbox')
        for shop in shop_list:
            try:
                shop_name = shop.select_one('.info-dealername a').text
                shop_id = re.search('[\d]+', urllib.parse.urlparse(shop.select_one('.info-dealername a').get('href')).path).group()
                super_shop = SuperShop(name=shop_name, shop_id=shop_id)
                super_shop_list.append(super_shop)
            except AttributeError as e:
                logger.error(f'action=scrape_shop_list_page error={e}')
                continue

        logger.info('action=scrape_shop_list_page status=done') 
        return super_shop_list

    @staticmethod
    def scrape_next_page_url(response: str) -> str|None:
        logger.info('action=scrape_next_page_url status=run')

        soup = BeautifulSoup(response, 'lxml')
        next_page_url = soup.select_one('.page-nav-next')
        if next_page_url:
            try:
                next_page_url = urllib.parse.urljoin(settings.SUPER_DOMAIN_URL, next_page_url.attrs.get('href'))
            except KeyError as e:
                logger.error(f'next_page_selector error={e}')
                next_page_url = None
        else:
            next_page_url = None

        logger.info('action=scrape_next_page_url status=run')
        return next_page_url

    @staticmethod
    def scrape_favorite_product_list_page(response: str) -> list[str, int]:
        logger.info('action=scraping_favorite_page status=run')

        products = []
        soup = BeautifulSoup(response, 'lxml')
        items = soup.select('.itembox-out-line')

        for item in items:
            try:
                url = settings.SUPER_DOMAIN_URL + item.select_one('.title a').attrs.get('href')
                cost = int(int(''.join(re.findall('\\d+', item.select_one('.trade-status-large').text))) * 1.1)
                products.append([url, cost])
            except AttributeError as e:
                logger.error(e)
                continue

        return products
