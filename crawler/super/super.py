import re
import time
import copy
from datetime import datetime
from datetime import timedelta
import logging.config

from requests import Session
from requests import Response
import requests
import openpyxl
from bs4 import BeautifulSoup

import settings
from crawler import utils
from crawler.super.models import Super
from crawler.super.models import SuperShop
from ims.models import FavoriteProduct


logger = logging.getLogger(__name__)


number_regex = re.compile('\\d+')

def login() -> Session:
    logger.info('action=login status=run')

    # session = requests.session()
    # response = utils.request(settings.LOGIN_URL, method='POST', session=session)
    
    # return session

    for _ in range(60):
        try:
            session = requests.session()
            response = session.post(settings.SUPER_LOGIN_URL, data=settings.SUPER_LOGIN_INFO)
            time.sleep(2)
            if not response.status_code == 200 or response is None:
                raise Exception
            return session
        except Exception as e:
            logger.error(f'action=login error={e}')
            time.sleep(30)


def super_main(shop_url, save_path=settings.SCRAPE_SAVE_PATH):
    logger.error('action=super_main status=run')

    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    session = login()
    products = get_url_list_page(session, shop_url)
    db_list, web_list = classify_exist_db(products)
    detail_list = get_product_detail_page(session, web_list)
    detail_list.extend(db_list)
    save_path = save_path + 'super' + timestamp + '.xlsx'
    list_to_excel_file(detail_list, save_path)


def get_url_list_page(session, shop_url):
    logger.info('action=get_url_list_page status=run')

    products = []

    while True:
        response = utils.request(url=shop_url, session=session)
        time.sleep(2)
        product_list = list_page_selector(response)
        next_url = next_page_selector(response)
        products.extend(product_list)
        if next_url is None:
            return products
        shop_url = next_url


def get_product_detail_page(session, products):
    logger.info('action=get_product_detail_page status=run')

    result_list = []
    for product in products:
        response = utils.request(url=product.url, session=session)
        time.sleep(2)
        product_dict = detail_page_selector(response)
        logger.debug(product_dict)
        for key, value in product_dict.items():
            copy_product = copy.deepcopy(product)
            copy_product.jan = key
            copy_product.price = value
            item = Super.get_product_jan_and_update_price(copy_product.product_code,
                                                          copy_product.jan, copy_product.price)
            if item is None:
                copy_product.save()
            result_list.append(copy_product)

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

    workbook = openpyxl.Workbook()
    sheet = workbook['Sheet']
    sheet.append(['JAN', 'Cost'])

    for product in products:
        if product.jan and product.price:
            sheet.append([product.jan, product.price])

    if sheet.max_row == 1:
        logger.info("This Shop don't have JAN_CODE")
        workbook.close()
        return

    workbook.save(save_path)
    workbook.close()


def detail_page_selector(response: Response):
    logger.info('action=detail_page_selector status=run')

    soup = BeautifulSoup(response.text, 'lxml')
    table = soup.select('.ts-tr02')

    products = {}

    for row in table:
        try:
            jan = ''.join(re.findall('[0-9]{13}', row.select_one('.co-fcgray.td-jan').text))
            price = int(int(''.join(re.findall('\\d+', row.select_one('.td-price02').text))) * 1.1)
        except AttributeError as e:
            logger.error(e)
            continue
        flg = products.get(jan)
        if flg is None or flg > price:
            products[jan] = price

    return products


def list_page_selector(response: Response):
    logger.info('action=list_page_selector status=run')

    result_list = []

    soup = BeautifulSoup(response.text, 'lxml')
    products = soup.select('.itembox-parts')
    for product in products:
        try:
            item_name = product.select_one('.item-name a')
            name = item_name.text.strip().replace('\u3000', '')
            url = settings.SUPER_DOMAIN_URL + item_name.attrs.get('href')
            product_code = item_name.attrs['href'].split('/')[-2]
            shop_code = response.url.split('/')[6]
            price = product.select_one('.item-price').text
            price = int(int(''.join(re.findall('\\d+', price))) * 1.1)
            item = Super.create(name=name, url=url, product_code=product_code, shop_code=shop_code, price=price)

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
            next_page_url = settings.SUPER_DOMAIN_URL + next_url.attrs['href']
        except KeyError as e:
            logger.error(f'next_page_selector error={e}')
            next_page_url = None
    else:
        next_page_url = None
    return next_page_url


def schedule_super_task():
    yesterday = datetime.now() - timedelta(days=1)
    url = settings.SUPER_NEW_PRODUCTS_URL + yesterday.strftime("%Y%m%d")
    super_main(url, save_path=settings.SCRAPE_SCHEDULE_SAVE_PATH)


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
