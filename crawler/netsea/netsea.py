import json
import time
import re
import os
import threading
from datetime import datetime
from queue import Queue
from urllib.parse import urljoin
from urllib.parse import urlparse
from urllib.parse import parse_qs
from typing import List
from collections import deque

import requests
from requests import Session
from bs4 import BeautifulSoup
import pandas as pd

import settings
import log_settings
from mq import MQ
from crawler.netsea.models import NetseaProduct
from crawler import utils
from crawler.netsea.models import NetseaShop


logger = log_settings.get_logger(__name__)

price_regex = re.compile('\\d+')
jan_regex = re.compile('[0-9]{13}')


class Netsea(object):

    def __init__(self, urls: List[str], timestamp: datetime = datetime.now(), is_new_product_search: bool = False):
        self.start_urls = deque(urls)
        self.url = ''
        self.netsea_product_queue = Queue()
        self.mq = MQ('mws')
        self.session = self.login()
        self.timestamp = timestamp.strftime("%Y%m%d_%H%M%S")
        self.is_new_product_search = is_new_product_search

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

    def pool_product_list_page(self, interval_sec: int = 2) -> None:
        logger.info('action=pool_product_list_page status=run')

        while self.url is not None or self.start_urls:
            if not self.url:
                self.url = self.start_urls.popleft()
            logger.info(self.url)
            response = utils.request(session=self.session, url=self.url)
            time.sleep(interval_sec)
            [self.netsea_product_queue.put(product) for product in (NetseaHTMLPage.scrape_product_list_page(response.text))]
            self.url = NetseaHTMLPage.scrape_next_page_url(response.text, response.url, self.is_new_product_search)
        
        self.netsea_product_queue.put(None)
        logger.info('action=pool_product_list_page status=done')

    def pool_product_detail_page(self, interval_sec: int = 2):
        logger.info('action=pool_product_detail_page status=run')

        while True:
            netsea_product = self.netsea_product_queue.get()
            if netsea_product is None:
                break
            if re.fullmatch('[0-9]{13}', netsea_product.product_code):
                netsea_product.jan = netsea_product.product_code
                netsea_product.save()
            else:
                product = NetseaProduct.get_object_filter_productcode_and_shopcode(netsea_product.product_code, netsea_product.shop_code)
                if product:
                    netsea_product.jan = product.jan
                else:
                    url = urljoin(settings.NETSEA_SHOP_URL, f'{netsea_product.shop_code}/{netsea_product.product_code}')
                    logger.info({'url': url})
                    response = utils.request(session=self.session, url=url)
                    time.sleep(interval_sec)
                    if response is None:
                        continue
                    parsed_value = NetseaHTMLPage.scrape_product_detail_page(response.text)
                    netsea_product.jan = parsed_value.get('jan')
                    netsea_product.save()
            
            self.publish_queue(netsea_product.jan, netsea_product.price)

        logger.info('action=pool_product_detail_page status=done')

    def pool_favorite_product_list_page(self, interval_sec: int = 2) -> pd.DataFrame:
        logger.info('action=pool_favorite_product_list_page status=run')

        netsea_product_list = []

        while self.url is not None:
            response = utils.request(session=self.session, url=self.url)
            time.sleep(interval_sec)
            netsea_product_list.extend(NetseaHTMLPage.scrape_favorite_list_page(response.text))
            self.url = NetseaHTMLPage.scrape_next_page_url(response.text, response.url)

        df = pd.DataFrame(data=None, columns={'jan': str, 'cost': int})
        PRODUCT_CODE_NUM = -1
        SHOP_ID_NUM = -2

        for product in netsea_product_list:
            url, price = product
            product_id = url.split('/')[PRODUCT_CODE_NUM]
            shop_id = url.split('/')[SHOP_ID_NUM]
            netsea_object = NetseaProduct.get_object_filter_productcode_and_shopcode(product_id, shop_id)
            if netsea_object and netsea_object.jan:
                jan = netsea_object.jan
            elif re.fullmatch('[0-9]{13}', product_id):
                jan = product_id
            else:
                continue
            df = df.append({'jan': jan, 'cost': price}, ignore_index=True)
        df = df.dropna()
        logger.info('action=pool_favorite_product_list_page status=done')
        return df

    def publish_queue(self, jan: str, price: int) -> None:
        logger.info('action=publish_queue status=run')

        if not jan or not price:
            return None

        params = {
                'filename': f'netsea_{self.timestamp}',
                'jan': jan,
                'cost': price,
            }
        self.mq.publish(json.dumps(params))
        logger.info('action=publish_queue status=done')

    def start_search_products(self):
        logger.info('action=start_search_products status=run')

        thread = threading.Thread(target=self.pool_product_detail_page, )
        thread.start()
        self.pool_product_list_page()
        thread.join()

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
            jan = re.fullmatch('[0-9]{13}', soup.select('#itemDetailSec td')[-1].text.strip())
            jan = jan.group()
        except (IndexError, AttributeError) as e:
            logger.error(f'action=get_jan error={e}')
            jan = None

        logger.info('action=scrape_detail_product_page status=done')
        return {'jan': jan}

    @classmethod
    def scrape_next_page_url(cls, response: str, response_url: str, is_new_product_search: bool = False) -> str | None:
        logger.info('action=scrape_next_page_url status=run')

        soup = BeautifulSoup(response, 'lxml')
        try:
            next_page_url_tag = soup.select_one('.next a')
            products = soup.select('.showcaseType01')
            new_product_count = soup.select('.showcaseHd .labelType04')
        except AttributeError as e:
            logger.error(f"action=next_page_url_selector status={e}")
            return None

        if is_new_product_search and (not len(new_product_count) == 60 or not next_page_url_tag):
            logger.info(f'next_page_url is None or new product flag is None')
            logger.info({'next_page_url': next_page_url_tag, 'new_product_count': len(new_product_count)})
            return None

        if next_page_url_tag:
            next_page_url = urljoin(settings.NETSEA_NEXT_URL, next_page_url_tag.attrs.get('href'))
        elif len(products) == 60:
            price = int(''.join(price_regex.findall(products[-1].select_one('.price').text)))
            current_url = urlparse(response_url)
            query = parse_qs(current_url.query)
            query['page'] = ['1']
            facet_price_to = query.get('facet_price_to')
            if facet_price_to == str(price):
                query['facet_price_to'] = str(int(facet_price_to.pop()) - 1)
            else:
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
