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
from crawler import netsea, utils
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
                self.publish_queue(netsea_product.jan, netsea_product.price, netsea_product.url)
                continue

            product = NetseaProduct.get_object_filter_productcode_and_shopcode(netsea_product.product_code, netsea_product.shop_code)
            if product:
                netsea_product.jan = product.jan
                product.url = netsea_product.url
                product.save()
                self.publish_queue(netsea_product.jan, netsea_product.price, netsea_product.url)
                continue

            logger.info({'url': netsea_product.url})
            response = utils.request(session=self.session, url=netsea_product.url)
            time.sleep(interval_sec)
            if response is None:
                continue

            parsed_value = NetseaHTMLPage.scrape_product_detail_page(response.text)
            netsea_product.jan = parsed_value.get('jan')
            netsea_product.save()

            self.publish_queue(netsea_product.jan, netsea_product.price, netsea_product.url)

        logger.info('action=pool_product_detail_page status=done')

    def start_favorite_products(self, url: str, interval_sec: int=2) -> None:
        logger.info({"action": "start_favorite_products", "status": "run"})

        products = []
        while url is not None:
            response = utils.request(session=self.session, url=url, time_sleep=interval_sec)
            products.extend(NetseaHTMLPage.scrape_favorite_list_page(response.text))
            url = NetseaHTMLPage.scrape_next_page_url(response.text, response.url)
            logger.info({"action": "start_favorite_products",
                         "messages": f"next_page_url is {url}"})

        for product in products:
            p = NetseaProduct.get_object_filter_productcode_and_shopcode(product.product_code, product.shop_code)
            if p is None:
                logger.error({"action": "start_favorite_products", 
                              "message": "Not Found product in database",
                              "product_code": product.product_code,
                              "shop_code": product.shop_code,})
                continue
            if not p.jan:
                logger.error({"action": "start_favorite_products", 
                              "message": "product has'nt jan code",
                              "product_code": product.product_code,
                              "shop_code": product.shop_code,})
                continue
            product.jan = p.jan

        [self.publish_queue(product.jan, product.price, product.url) for product in products] 

        logger.info({"action": "start_favorite_products", "status": "done"})

    def publish_queue(self, jan: str, price: int, url: str) -> None:
        logger.info('action=publish_queue status=run')

        if not all([jan, price, url]):
            return

        self.mq.publish(json.dumps({
                'filename': f'netsea_{self.timestamp}',
                'jan': jan,
                'cost': price,
                'url': url,
            }))

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

            url = product.select_one('.showcaseHd a').attrs.get('href')
            if url:
                url = urljoin(settings.NETSEA_ENDPOINT, url)
            shop_code = urlparse(url).path.split('/')[SHOP_CODE_NUM]
            product_code = urlparse(url).path.split('/')[PRODUCT_CODE_NUM]
            netsea_product = NetseaProduct(
                                        name=title, 
                                        price=price,
                                        shop_code=shop_code,
                                        product_code=product_code,
                                        url=url,
                                    )

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
            price = None
            for product in products[-1:]:
                price_str = product.select_one('.price')
                if price_str:
                    price = int(''.join(price_regex.findall(price_str.text)))
                    break
            if price is None:
                return None

            current_url = urlparse(response_url)
            query = parse_qs(current_url.query)
            query['page'] = ['1']
            facet_price_to = query.get('facet_price_to')
            if facet_price_to == str(price):
                query['facet_price_to'] = str(int(facet_price_to.pop()) - 1)
            else:
                query['facet_price_to'] = str(price - 1)
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
    def scrape_favorite_list_page(cls, response: str, tax_rate: float=1.1) -> list[NetseaProduct]:
        logger.info('action=scraping_favorite_list_page status=run')

        # e.g. https://www.netsea.jp/shop/84918/28234210
        SHOP_CODE_INDEX = 1
        PRODUCT_CODE_INDEX = 2

        products = []
        soup = BeautifulSoup(response, 'lxml')
        products_box = soup.select('form .showcaseType03')

        for box in products_box:
            title_tag = box.select_one('.showcaseHd a')
            if not title_tag:
                logger.error({"action": "scrape_favorite_list_page", "message": "Not Found Title"})
                continue
            title = title_tag.text.strip()
            url = title_tag.attrs.get("href")
            if not url:
                logger.error({"action": "scrape_favorite_list_page", "message": "Not Found URL"})
                continue

            url_path = list(filter(None, urlparse(url).path.split("/")))
            try:
                shop_code = url_path[SHOP_CODE_INDEX]
                product_code = url_path[PRODUCT_CODE_INDEX]
            except IndexError as ex:
                logger.error({
                    "action": "scrape_favorite_list_page",
                    "message": "Not Found product_code and shop_code",
                    "error": ex})
                continue

            price = price_tag if (price_tag := box.select_one('.afterPrice')) else box.select_one(".price")
            if price:
                price = int(int(''.join(re.findall('[0-9]+', price.text))) * tax_rate)

            products.append(NetseaProduct(name=title, price=price,
                            shop_code=shop_code, product_code=product_code, url=url))
        logger.info('action=scraping_favorite_list_page status=done')
        return products
