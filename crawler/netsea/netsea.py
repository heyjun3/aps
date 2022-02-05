import time
import re
import os
import logging
import datetime
import logging.config
from urllib.parse import urljoin

import requests
from requests import Session
from bs4 import BeautifulSoup
import openpyxl
from requests import Response
import pandas as pd

import settings
from crawler.netsea.models import Netsea
from crawler import utils
from crawler.netsea.models import NetseaShop
from ims.models import FavoriteProduct


logger = logging.getLogger(__name__)

price_regex = re.compile('\\d+')
jan_regex = re.compile('[0-9]{13}')


def login() -> Session:
    try:
        session = requests.session()
        response = session.get(settings.NETSEA_LOGIN_URL)
        soup = BeautifulSoup(response.text, 'html.parser')

        authenticity = soup.find(attrs={'name': '_token'}).get('value')
        cookie = response.cookies
        info = {
            '_token': authenticity,
            'login_id': settings.NETSEA_ID,
            'password': settings.NETSEA_PASSWD,
        }
        session.post(settings.NETSEA_LOGIN_URL, data=info, cookies=cookie)
        time.sleep(2)
        return session
    except Exception as e:
        logger.error(f'action=login error={e}')
        raise


def list_page_selector(response, new_bool):
    logger.info('action=list_page_selector status=run')

    result_list = []
    new_flag = True

    soup = BeautifulSoup(response.content, 'lxml')
    product_list = soup.select('.showcaseType01')

    for product in product_list:
        url = product.select_one('.showcaseHd a')
        price = product.select_one('.price')
        flag = product.select_one('.labelType04')

        if new_bool and flag is None:
            new_flag = False
            return result_list, new_flag

        if url is None or price is None:
            logger.info('url_title or price is None')
            continue

        item = Netsea()
        item.name = url.text.strip()
        item.price = int(int(''.join(price_regex.findall(price.text))) * 1.1)
        item.url = url.attrs['href'].split('?')[0]
        item.shop_code = item.url.split('/')[-2]
        item.product_code = item.url.split('/')[-1]
        result_list.append(item)

    return result_list, new_flag


def next_page_url_selector(response):
    logger.info('action=next_page_url_selector status=run')

    soup = BeautifulSoup(response.content, 'lxml')
    try:
        next_page_url = soup.select_one('.next a')
        current = soup.select_one('.current').text.strip()
        products = soup.select('.showcaseType01')
        is_beginner_trade_false = soup.select('.priceHide')
    except AttributeError as e:
        logger.error(f"action=next_page_url_selector status={e}")
        return None

    if len(products) == len(is_beginner_trade_false):
        return None

    if next_page_url:
        next_page_url = urljoin(settings.NETSEA_NEXT_URL, next_page_url.attrs['href'])

    if current == '166' and len(products) == 60:
        price = products[-1].select_one('.price')
        price = ''.join(price_regex.findall(price.text))
        try:
            supplier_id = re.findall('supplier_id=[\\d]+', response.url)[0]
            supplier_id = ''.join(price_regex.findall(supplier_id))
            next_page_url = urljoin(settings.NETSEA_NEXT_URL, f"?supplier_id={supplier_id}&sort=PD&facet_price_to={price}")
        except IndexError as e:
            logging.error(f'action=next_page_selector status={e}')
            return None

    return next_page_url


def get_url_cost_list_page(session: Session, url: str, new_bool=False) -> list:
    logger.info('action=get_url_cost_list_page status=run')

    result_list = []
    while True:
        response = utils.request(session=session, url=url)
        time.sleep(2)
        products, new_flag = list_page_selector(response, new_bool)
        next_page_url = next_page_url_selector(response)
        result_list.extend(products)

        if not new_flag:
            logger.info('new products end')
            break

        if next_page_url is None:
            logger.info('next_page_url is None. break all_url_cost_list')
            break
        url = next_page_url

    return result_list


def classify_exist_jan_url(products: list):
    logger.info('action=classify_exist_jan_url status=run')

    not_exist_jan, exist_jan = [], []

    for product in products:
        jan = product.url.split('/')[-1]
        jan_flag = re.fullmatch('[0-9]{13}', jan)
        if jan_flag is None:
            not_exist_jan.append(product)
        else:
            product.jan = jan
            exist_jan.append(product)

    return not_exist_jan, exist_jan


def classify_exist_db(products: list):
    logger.info('action=classify_exist_db status=run')
    logger.debug(products)

    if not products:
        logger.info('args is False')
        return [], []

    not_exist_list, exist_list = [], []

    for product in products:
        db_response = Netsea.get(product.url)
        if not db_response:
            not_exist_list.append(product)
        else:
            exist_list.append(db_response)
    return exist_list, not_exist_list


def detail_page_selector(response: Response):
    logger.info('action=detail_page_selector status=run')
    soup = BeautifulSoup(response.content, 'lxml')
    try:
        jan = soup.select('#itemDetailSec td')[-1]
    except IndexError as e:
        logger.error(f'action=get_jan error={e}')
        return None

    jan = ''.join(jan_regex.findall(jan.text))
    logger.debug(jan)
    return jan


def get_detail_page_jan(session, products):
    logger.info('action=get_jan status=run')

    if not products:
        return products

    for product in products:
        response = utils.request(session=session, url=product.url)
        time.sleep(2)
        product.jan = detail_page_selector(response)
        product.save()

    return products


def list_to_excel_file(products: list, save_path: str) -> None:
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


def shop_list_page_selector(response: Response):
    logger.info('action=shop_list_page_selector status=run')
    new_shop_urls = []

    soup = BeautifulSoup(response.text, 'html.parser')
    shops = soup.select('.supNameList a')
    category = re.search('[0-9]', response.url).group()

    for shop in shops:
        shop_name = shop.text
        shop_url = shop.attrs['href']
        shop_id = int(shop_url.split('/')[-1])
        create_bool = NetseaShop.create(name=shop_name, shop_id=shop_id, url=shop_url,
                                        quantity=None, category=category)
        if create_bool:
            new_shop_urls.append(shop_url)

    next_url = response.url.replace(category, str(int(category)+1))
    if category == '9':
        next_url = None

    return new_shop_urls, next_url


# def new_shop_search():
#     logger.info('action=new_shop_search status=run')
#     session = requests.session()
#     start_url = 'https://www.netsea.jp/shop?category_id=1&sort=NEW'
#     new_shop_urls = common.common_pool_list_page(shop_list_page_selector, session, start_url)
#     for url in new_shop_urls:
#         run_netsea(url)


# def get_product_count(session: Session, url: str):
#     logger.info('action=get_product_count status=run')
#     response = session_request(session=session, url=url)
#     time.sleep(2)
#     soup = BeautifulSoup(response.text, 'html.parser')
#     product_quantity = soup.select_one('.currentCate span')
#     if product_quantity is None:
#         return None
#     product_quantity = ''.join(price_regex.findall(product_quantity.text))
#     return product_quantity

def new_product_search():
    logger.info('action=new_product_search status=run')

    result_list = []
    timestamp = datetime.datetime.now().strftime("%Y%m%d_%H%M%S")
    session = login()

    for i in range(2, 8):
        url = f'https://www.netsea.jp/search?sort=new&category_id={i}'
        logger.debug(url)
        products_list = get_url_cost_list_page(session=session, url=url, new_bool=True)
        object_not_have_jan, url_include_jan = classify_exist_jan_url(products_list)

        db_list, web_list = classify_exist_db(object_not_have_jan)

        web_list = get_detail_page_jan(session, web_list)
        url_include_jan.extend(web_list)
        result_list.extend(url_include_jan)

    save_path = os.path.join(settings.SCRAPE_SCHEDULE_SAVE_PATH, f'netsea{timestamp}.xlsx')

    list_to_excel_file(result_list, save_path)


def run_netsea(url: str, discount_rate=1.0):
    logger.info('action=run_netsea status=run')

    timestamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')
    session = login()
    products_list = get_url_cost_list_page(session=session, url=url+"?sort=PD")
    object_not_have_jan, url_include_jan = classify_exist_jan_url(products_list)

    db_list, web_list = classify_exist_db(object_not_have_jan)
    url_include_jan.extend(db_list)
    logger.info('classify_exist_db is done')

    web_list = get_detail_page_jan(session, web_list)
    url_include_jan.extend(web_list)

    if not discount_rate == 1.0 and discount_rate is not None:
        rate = [discount_rate for _ in range(len(url_include_jan))]
        url_include_jan = list(map(calc_discount_price, url_include_jan, rate))

    save_path = os.path.join(settings.SCRAPE_SCHEDULE_SAVE_PATH, f'netsea{timestamp}.xlsx')

    list_to_excel_file(url_include_jan, save_path)


def shop_product_quantity_selector(response):
    soup = BeautifulSoup(response.content, 'lxml')
    product_quantity = soup.select_one('.currentCate span')
    if product_quantity is None:
        return None
    product_quantity = ''.join(price_regex.findall(product_quantity.text))
    return int(product_quantity)


def calc_discount_price(product: Netsea, discount_rate: float):
    product.price = int(int(product.price) * discount_rate)
    return product


def run_discount():
    print('Enter excel path')
    path = input()
    book = openpyxl.load_workbook(path)
    sheet = book[book.sheetnames[0]]
    for row in sheet.values:
        if type(row[0]) == int:
            run_netsea(settings.NETSEA_SHOP_URL+str(row[0]), discount_rate=0.95)


def discount_shops():
    logger.info('action=discount_shops status=run')
    path = os.path.join(settings.BASE_PATH, 'QBizeLc6MSjZ.xlsx')
    book = openpyxl.load_workbook(path)
    sheet = book[book.sheetnames[0]]
    shop_ids = []

    for row in sheet.values:
        if type(row[0]) == int:
            shop_ids.append(row[0])

    return shop_ids


def collect_favorite_products(interval_sec: int = 2):
    collect_data = []
    url = 'https://www.netsea.jp/bookmark?stock_option=in&page=1'
    session = login()
    while True:
        response = utils.request(session=session, url=url)
        time.sleep(interval_sec)
        products = scraping_favorite_list_page(response)
        url = scraping_next_url_favorite_list_page(response)
        collect_data.extend(products)
        if not url:
            break

    df = pd.DataFrame(data=None, columns={'jan': str, 'cost': int})
    for product in collect_data:
        jan = ''
        url, price = product
        product_id = url.split('/')[-1]
        if re.fullmatch('[0-9]{13}', product_id):
            jan = product_id
        else:
            netsea_object = Netsea.get(url)
            if netsea_object is None:
                print(product)
                continue
            elif not netsea_object.jan:
                print(netsea_object.value)
            else:
                jan = netsea_object.jan
        df = df.append({'jan': jan, 'cost': price}, ignore_index=True)
    df = df.dropna()
    return df


def scraping_favorite_list_page(response: Response):
    logger.info('action=scraping_favorite_list_page status=run')

    product_data = []
    soup = BeautifulSoup(response.content, 'lxml')
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


def scraping_next_url_favorite_list_page(response: Response):
    logger.info('action=scraping_next_url_favorite_list_page status=run')

    soup = BeautifulSoup(response.content, 'lxml')
    try:
        next_url = soup.select_one('.next a').attrs.get('href')
    except AttributeError as e:
        logger.error(f'action=scraping_next_url_favorite_list_page error={e}')
        next_url = None

    logger.info('action=scraping_next_url_favorite_list_page status=done')
    return next_url
