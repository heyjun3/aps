import re
import time
import copy
import os
import urllib
from datetime import datetime
from datetime import timedelta
import logging.config

from requests import Session
from requests import Response
import requests
from bs4 import BeautifulSoup
import pandas as pd

import settings
from crawler import utils
from crawler.super.models import Super
from crawler.super.models import SuperShop
from crawler.super.models import SuperProductDetails
from ims.models import FavoriteProduct


logger = logging.getLogger(__name__)


number_regex = re.compile('\\d+')


def super_main(shop_url, save_path=settings.SCRAPE_SAVE_PATH):
    logger.info('action=super_main status=run')

    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    session = login()
    products = get_url_list_page(session, shop_url)
    products_detail_list = get_product_detail_page(session, products)
    save_path = os.path.join(save_path, f'super{timestamp}.xlsx')
    list_to_excel_file(products_detail_list, save_path)


def schedule_super_task():
    logger.info('action=schedule_super_task status=run')

    yesterday = datetime.now() - timedelta(days=1)
    url = settings.SUPER_NEW_PRODUCTS_URL + yesterday.strftime("%Y%m%d")
    super_main(url, save_path=settings.SCRAPE_SCHEDULE_SAVE_PATH)

    logger.info('action=schedule_super_task status=done')


def super_all():
    logger.info('action=super_all status=run')

    shops = SuperShop.get_all_info()
    for shop in shops:
        super_main(shop.url)

    logger.info('action=super_all status=done')


def login() -> Session:
    logger.info('action=login status=run')

    session = requests.session()
    response = utils.request(url=settings.SUPER_LOGIN_URL, method='POST', session=session, data=settings.SUPER_LOGIN_INFO)
    
    return session


def get_url_list_page(session, shop_url, interval_sec: int = 2):
    logger.info('action=get_url_list_page status=run')

    products = []
    next_url = shop_url

    while True:
        response = utils.request(url=next_url, session=session)
        time.sleep(interval_sec)
        product_list = list_page_selector(response)
        next_url = next_page_selector(response)
        products.extend(product_list)
        if next_url is None:
            logger.info('action=get_url_list_page status=done')
            return products


def get_product_detail_page(session, products, interval_sec: int = 2):
    logger.info('action=get_product_detail_page status=run')

    result_list = []
    for product in products:
        db_response = SuperProductDetails.get(product.product_code, product.price)
        if db_response is None:
            response = utils.request(url=product.url, session=session)
            time.sleep(interval_sec)
            product_list = detail_page_selector(response)
            result_list.extend(product_list)
        else:
            result_list.extend(db_response)

    return result_list


def classify_exist_db(products):
    logger.info('action=classify_exist_db status=run')

    exist_list, not_exist_list = [], []
    for product in products:
        db_response = Super.get_product(product.product_code, product.price)
        if not db_response:
            not_exist_list.append(product)
        else:
            exist_list.extend(db_response)

    return exist_list, not_exist_list


def list_to_excel_file(products: list, save_path: str):
    logger.info('action=list_to_excel_file status=run')

    jan_price_list = [[product.jan, product.price] for product in products if product.jan]
    df = pd.DataFrame(data=jan_price_list, columns=['JAN', 'Cost']).astype({'JAN': str, 'Cost': int}).drop_duplicates()

    if not df.empty:
        df.to_excel(save_path, index=False)    

    logger.info('action=list_to_excel_file status=done')


def detail_page_selector(response: Response, sales_tax: float = 1.1):
    logger.info('action=detail_page_selector status=run')
    products = []

    soup = BeautifulSoup(response.text, 'lxml')
    table = soup.select('.ts-tr02')
    product_code = ''.join(re.findall('\\d+', soup.select_one('.co-fs12.co-clf.reduce-tax .co-pc-only').text))
    shop_code = ''.join(re.findall('\\d+', soup.select_one('.dl-name-txt').get('href')))

    for row in table:
        try:
            jan = ''.join(re.findall('[0-9]{13}', row.select_one('.co-fcgray.td-jan').text))
        except AttributeError as e:
            logger.error(f"product hasn't jan code error={e}")
            jan = None
        price = int(int(''.join(re.findall('\\d+', row.select_one('.td-price02').text))) * sales_tax)
        set_number = int(re.search('\\d+', row.select_one('.co-align-center.co-pc-only.border-rt.border-b').text.strip()).group())
        super_product = SuperProductDetails(jan=jan, price=price, set_number=set_number)
        products.append(super_product)

    products = {product.jan: product for product in sorted(products, key=lambda x: x.price, reverse=True)}
    for product in products.values():
        product.product_code = product_code
        product.shop_code = shop_code
        product.save_or_update()
    
    products = list(products.values())
    
    return products

def list_page_selector(response: Response, sales_tax: float = 1.1):
    logger.info('action=list_page_selector status=run')

    result_list = []

    soup = BeautifulSoup(response.text, 'lxml')
    products = soup.select('.itembox-parts')
    for product in products:
        try:
            item_name = product.select_one('.item-name a')
            name = item_name.text.strip().replace('\u3000', '')
            url = urllib.parse.urljoin(settings.SUPER_DOMAIN_URL, item_name.attrs.get('href'))
            product_code = item_name.attrs['href'].split('/')[-2]
            shop_code = response.url.split('/')[6]
            price = product.select_one('.item-price').text
            price = int(int(''.join(re.findall('\\d+', price))) * sales_tax)
            item = Super(name=name, url=url, product_code=product_code, shop_code=shop_code, price=price)
            item.save()

        except AttributeError as e:
            logger.error(f'action=list_page_selector error={e}')
            continue
        result_list.append(item)

    return result_list


def shop_list_page_selector(response: Response):
    logger.info('action=shop_name_selector status=run')

    soup = BeautifulSoup(response.text, 'lxml')
    shop_list = soup.select('.dealer-eachbox')
    logger.debug(len(shop_list))
    for shop in shop_list:
        try:
            shop_name = shop.select_one('.info-dealername a').text
            shop_id = int(shop.select_one('.info-dealername a').get('href').split('/')[-2])
            shop_url = rf'https://www.superdelivery.com/p/do/dpsl/{shop_id}/'
            shop_quantity = ''.join(number_regex.findall(shop.select_one('.info-dealeritemnum a').text))
            SuperShop.create(name=shop_name, shop_id=shop_id, url=shop_url, quantity=shop_quantity, category=None)
            logger.debug(f'{shop_name}:{shop_id}:{shop_url}:{shop_quantity}')
        except AttributeError as e:
            logger.error(rf'action=shop_name_selector error={e}')
            continue

    next_page_url = next_page_selector(response)
    logger.debug(next_page_url)

    return shop_list, next_page_url


def next_page_selector(response):
    logger.info('action=next_page_selector status=run')

    soup = BeautifulSoup(response.content, 'lxml')
    next_url = soup.select_one('.page-nav-next')
    if next_url:
        try:
            next_page_url = urllib.parse.urljoin(settings.SUPER_DOMAIN_URL, next_url.attrs['href'])
        except KeyError as e:
            logger.error(f'next_page_selector error={e}')
            next_page_url = None
    else:
        next_page_url = None
    return next_page_url


def collection_favorite_products():

    collect_data = []
    url = 'https://www.superdelivery.com/p/wishlist/search.do'
    session = login()

    while True:
        response = utils.request(url=url, session=session)
        time.sleep(2)
        products = scraping_favorite_page(response)
        url = next_page_selector(response)
        collect_data.extend(products)
        logger.debug(url)
        if url is None:
            break

    for item in collect_data:
        url, cost = item
        products = Super.get_url(url)
        if products is None:
            continue
        else:
            for product in products:
                jan = product.jan
                if not jan:
                    print(product.value)
                FavoriteProduct.save(url=url, jan=jan, cost=cost)
    return None


def scraping_favorite_page(response: Response):
    logger.info('action=scraping_favorite_page status=run')

    products = []
    soup = BeautifulSoup(response.text, 'lxml')
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
