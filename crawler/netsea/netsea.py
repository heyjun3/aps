import json
from operator import ne
import time
import re
import os
from datetime import datetime
from urllib.parse import urljoin
from urllib.parse import urlparse
from urllib.parse import parse_qs

import requests
from requests import Session
from bs4 import BeautifulSoup
import openpyxl
from requests import Response
import pandas as pd

import settings
import log_settings
from mq import MQ
from crawler.netsea.models import NetseaProduct
from crawler import netsea, utils
from crawler.netsea.models import NetseaShop


logger = log_settings.get_logger(__name__)

price_regex = re.compile('\\d+')
jan_regex = re.compile('[0-9]{13}')


class Netsea(object):

    def __init__(self, url, params: dict = None, timestamp: datetime = datetime.now()):
        self.url = requests.Request(method='GET', url=url, params=params).prepare().url
        self.netsea_product_list = []
        self.mq = MQ('mws')
        self.session = self.login()
        self.timestamp = timestamp.strftime("%Y%m%d_%H%M%S")

    def get_authentication_token(self, session: requests.Session) -> str:
        logger.info('action=get_authentication_token status=run')
        
        response = utils.request(url=settings.NETSEA_LOGIN_URL, session=session)
        soup = BeautifulSoup(response.text, 'lxml')
        authenticity_token = soup.find(attrs={'name': '_token'}).get('value')

        logger.info('action=get_authentication_token status=done')
        return authenticity_token

    def login(self) -> Session:
        logger.info('action=login status=run')

        session = requests.Session()
        token = self.get_authentication_token(session)
        info = {
            '_token': token,
            'login_id': settings.NETSEA_ID,
            'password': settings.NETSEA_PASSWD,
        }
        response = utils.request(url=settings.NETSEA_LOGIN_URL, method='POST', session=session, data=info)
        time.sleep(2)

        logger.info('action=login status=done')
        return session

    def pool_product_list_page(self, is_new_product_search: bool = False, interval_sec: int = 2) -> None:
        logger.info('action=pool_product_list_page status=run')

        while self.url is not None:
            response = utils.request(session=self.session, url=self.url)
            time.sleep(interval_sec)
            self.netsea_product_list.extend(NetseaHTMLPage.scrape_product_list_page(response.text))
            self.url = NetseaHTMLPage.scrape_next_page_url(response.text, response.url, is_new_product_search)
        
        logger.info('action=pool_product_list_page status=done')

    def pool_product_detail_page(self, interval_sec: int = 2):
        logger.info('action=pool_product_detail_page status=run')

        for netsea_product in self.netsea_product_list:
            product = NetseaProduct.get_object_filter_productcode_and_shopcode(netsea_product.product_code, netsea_product.shop_code)
            if product:
                netsea_product.jan = product.jan
            elif re.fullmatch('[\d]{13}', netsea_product.product_code):
                netsea_product.jan = netsea_product.product_code
                netsea_product.save()
            else:
                url = urljoin(settings.NETSEA_SHOP_URL, f'{netsea_product.shop_code}/{netsea_product.product_code}')
                response = utils.request(session=self.session, url=url)
                time.sleep(interval_sec)
                if response is None:
                    continue
                netsea_product.jan = NetseaHTMLPage.scrape_product_detail_page(response.text)
                netsea_product.save()
            
            if netsea_product.jan is None:
                continue

            params = {
                'filename': f'netsea_{self.timestamp}',
                'jan': netsea_product.jan,
                'cost': netsea_product.price
            }
            self.mq.publish(json.dumps(params))

        logger.info('action=pool_product_detail_page status=done')

    def pool_favorite_product_list_page(self, interval_sec: int = 2) -> pd.DataFrame:
        logger.info('action=pool_favorite_product_list_page status=run')
        while self.url is not None:
            response = utils.request(session=self.session, url=self.url)
            time.sleep(interval_sec)
            self.netsea_product_list.extend(NetseaHTMLPage.scrape_favorite_list_page(response.text))
            self.url = NetseaHTMLPage.scrape_next_page_url(response.text, response.url)

        df = pd.DataFrame(data=None, columns={'jan': str, 'cost': int})
        PRODUCT_CODE_NUM = -1
        SHOP_ID_NUM = -2

        for product in self.netsea_product_list:
            url, price = product
            product_id = url.split('/')[PRODUCT_CODE_NUM]
            shop_id = url.split('/')[SHOP_ID_NUM]
            netsea_object = NetseaProduct.get_object_filter_productcode_and_shopcode(product_id, shop_id)
            if netsea_object and netsea_object.jan:
                jan = netsea_object.jan
            elif re.fullmatch('[\d]{13}', product_id):
                jan = product_id
            else:
                continue
            df = df.append({'jan': jan, 'cost': price}, ignore_index=True)
        df = df.dropna()
        logger.info('action=pool_favorite_product_list_page status=done')
        return df


    def start_search_products(self):
        logger.info('action=start_search_products status=run')

        self.pool_product_list_page()
        self.pool_product_detail_page()

        logger.info('action=start_search_products status=done')


class NetseaHTMLPage(object):

    @classmethod
    def scrape_product_list_page(cls, response: str, consume_tax_rate: float = 1.1) -> list[NetseaProduct]:
        logger.info('action=scrape_product_list_page status=run')

        netsea_product_list = []
        SHOP_CODE_NUM = -2
        PRODUCT_CODE_NUM = -1

        soup = BeautifulSoup(response, 'lxml')
        product_list = soup.select('.showcaseType01')

        for product in product_list:
            try:
                title = product.select_one('.showcaseHd a').text.strip()
            except AttributeError as ex:
                logger.error(f'title is None error={ex}')
                continue

            try:
                price = product.select_one('.afterPrice')
                if price is None:
                    price = product.select_one('.price')
                price = int(int(''.join(price_regex.findall(price.text))) * consume_tax_rate)
            except AttributeError as ex:
                logger.error('price is None')
                continue

            url = urlparse(product.select_one('.showcaseHd a').attrs.get('href'))
            shop_code = url.path.split('/')[SHOP_CODE_NUM]
            product_code = url.path.split('/')[PRODUCT_CODE_NUM]
            netsea_product = NetseaProduct(name=title, price=price, shop_code=shop_code, product_code=product_code)

            netsea_product_list.append(netsea_product)
        
        logger.info('action=scrape_product_list_page status=done')
        return netsea_product_list

    @classmethod
    def scrape_product_detail_page(self, response: str) -> str | None:
        logger.info('action=scrape_detail_product_page status=run')
        soup = BeautifulSoup(response, 'lxml')
        try:
            jan = soup.select('#itemDetailSec td')[-1]
        except IndexError as e:
            logger.error(f'action=get_jan error={e}')
            return None

        jan = ''.join(jan_regex.findall(jan.text))
        logger.info('action=scrape_detail_product_page status=done')
        return jan

    @classmethod
    def scrape_next_page_url(cls, response: str, response_url: str, is_new_product_search: bool = False, consume_tax: float = 1.1) -> str | None:
        logger.info('action=scrape_next_page_url status=run')

        soup = BeautifulSoup(response, 'lxml')
        try:
            next_page_url_tag = soup.select_one('.next a')
            products = soup.select('.showcaseType01')
            new_product_count = soup.select('.labelType04')
        except AttributeError as e:
            logger.error(f"action=next_page_url_selector status={e}")
            return None

        if is_new_product_search and (not len(new_product_count) == 60 or not next_page_url_tag):
            return None

        if next_page_url_tag:
            next_page_url = urljoin(settings.NETSEA_NEXT_URL, next_page_url_tag.attrs.get('href'))
        elif len(products) == 60:
            price = int(int(''.join(price_regex.findall(products[-1].select_one('.price').text))) * consume_tax)
            current_url = urlparse(response_url)
            query = parse_qs(current_url.query)
            query['page'] = ['1']
            query['facet_price_to'] = str(price)
            next_page_url = requests.Request(url=settings.NETSEA_NEXT_URL, params=query).prepare().url
        else:
            next_page_url = None
        
        logger.info('action=scrape_next_page_url status=done')
        return next_page_url

    @classmethod
    def scrape_shop_list_page(cls, response: str) -> list[NetseaShop]:
        logger.info('action=shop_list_page_selector status=run')

        SHOP_ID_NUM = -1
        shop_list = []
        soup = BeautifulSoup(response, 'lxml')
        shops = soup.select('.supNameList a')

        for shop in shops:
            shop_name = shop.text
            shop_url = shop.attrs.get('href')
            shop_id = os.path.split(urlparse(shop_url).path)[SHOP_ID_NUM]
            netsea_shop = NetseaShop(name=shop_name, shop_id=shop_id)
            shop_list.append(netsea_shop)
        
        return shop_list

    @classmethod
    def scrape_favorite_list_page(cls, response: str) -> list[str, int]:
        logger.info('action=scraping_favorite_list_page status=run')

        product_data = []
        soup = BeautifulSoup(response, 'lxml')
        products_box = soup.select('form .showcaseType03')

        for box in products_box:
            try:
                url = box.select_one('.showcaseHd a').attrs.get('href')
                price = box.select_one('.afterPrice')
                if price is None:
                    price = box.select_one('.price')
                price = int(int(''.join(re.findall('[\\d+]', price.text))) * 1.1)
                product_data.append([url, price])
            except AttributeError as e:
                logger.error(f'action=scraping_favorite_list_page error={e}')

        logger.info('action=scraping_favorite_list_page status=done')
        return product_data 
