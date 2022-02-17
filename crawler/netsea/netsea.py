import time
import re
import os
import datetime
from urllib.parse import urljoin
from urllib.parse import urlparse
from urllib.parse import urlencode
from urllib.parse import parse_qs
from xml.dom.minidom import Attr

import requests
from requests import Session
from bs4 import BeautifulSoup
import openpyxl
from requests import Response
import pandas as pd

import settings
import log_settings
from crawler.netsea.models import NetseaProduct
from crawler import utils
from crawler.netsea.models import NetseaShop


logger = log_settings.get_logger(__name__)

price_regex = re.compile('\\d+')
jan_regex = re.compile('[0-9]{13}')


class Netsea(object):

    def __init__(self, url):
        self.url = url
        self.netsea_product_list = []
        self.next_page_url = ''

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

    def scrape_product_list_page(self, response: requests.Response) -> None:
        logger.info('action=list_page_selector status=run')

        SHOP_CODE_NUM = -2
        PRODUCT_CODE_NUM = -1

        soup = BeautifulSoup(response.content, 'lxml')
        product_list = soup.select('.showcaseType01')

        for product in product_list:
            try:
                title = product.select_one('.showcaseHd a').text.strip()
            except AttributeError as ex:
                logger.error('title is None')
                continue

            try:
                price = int(int(''.join(price_regex.findall(product.select_one('.price').text))) * 1.1)
            except AttributeError as ex:
                logger.error('price is None')
                continue

            url = urlparse(product.select_one('.showcaseHd a').attrs.get('href'))
            shop_code = url.path.split('/')[SHOP_CODE_NUM]
            product_code = url.path.split('/')[PRODUCT_CODE_NUM]
            netsea_product = NetseaProduct(name=title, price=price, shop_code=shop_code, product_code=product_code)

            self.netsea_product_list.append(netsea_product)

    def scrape_next_page_url(self, response: requests.Response, is_new_product_search: bool = False) -> None:
        logger.info('action=next_page_url_selector status=run')

        soup = BeautifulSoup(response.content, 'lxml')
        try:
            next_page_url_tag = soup.select_one('.next a')
            products = soup.select('.showcaseType01')
            new_product_count = soup.select('.labelType04')
        except AttributeError as e:
            logger.error(f"action=next_page_url_selector status={e}")
            self.next_page_url = None

        if is_new_product_search and not len(new_product_count) == 60:
            self.next_page_url = None

        if next_page_url_tag:
            self.next_page_url = urljoin(settings.NETSEA_NEXT_URL, next_page_url_tag.attrs.get('href'))
        elif len(products) == 60:
            price = ''.join(price_regex.findall(products[-1].select_one('.price').text))
            current_url = urlparse(response.url)
            query = parse_qs(current_url.query)
            query['page'] = ['1']
            query['facet_price_to'] = price
            self.next_page_url = requests.Request(url=settings.NETSEA_NEXT_URL, params=query).prepare().url

    def pool_product_list_page(self, session: Session, url: str, is_new_product_search: bool = False, interval_sec: int = 2) -> None:
        logger.info('action=get_url_cost_list_page status=run')

        while True:
            response = utils.request(session=session, url=url)
            time.sleep(interval_sec)
            self.scrape_product_list_page(response)
            self.scrape_next_page_url(response, is_new_product_search)

            if self.next_page_url is None:
                logger.info('next_page_url is None. break all_url_cost_list')
                break
            url = self.next_page_url

    def scrape_product_detail_page(self, response: Response) -> str:
        logger.info('action=scrape_detail_product_page status=run')
        soup = BeautifulSoup(response.content, 'lxml')
        try:
            jan = soup.select('#itemDetailSec td')[-1]
        except IndexError as e:
            logger.error(f'action=get_jan error={e}')
            return None

        jan = ''.join(jan_regex.findall(jan.text))
        logger.info('action=scrape_detail_product_page status=done')
        return jan


    def pool_product_detail_page(self, session: requests.Session, interval_sec: int = 2):
        logger.info('action=get_jan status=run')

        for netsea_product in self.netsea_product_list:
            response = utils.request(session=session, url=product.url)
            time.sleep(interval_sec)
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

    soup = BeautifulSoup(response.text, 'html.parser')
    shops = soup.select('.supNameList a')
    category = re.search('[0-9]', response.url).group()

    for shop in shops:
        shop_name = shop.text
        shop_url = shop.attrs['href']
        shop_id = int(shop_url.split('/')[-1])
        NetseaShop.create(name=shop_name, shop_id=shop_id, url=shop_url, quantity=None, category=category)
    
    next_url = response.url.replace(category, str(int(category)+1))
    if category == '9':
        return None

    return next_url


def new_shop_search():
    logger.info('action=new_shop_search status=run')
    session = login()
    url = 'https://www.netsea.jp/shop?category_id=1&sort=NEW'
    while True:
        response = utils.request(url=url, session=session)
        url = shop_list_page_selector(response)
        if url is None:
            logger.info('action=new_shop_search status=done')
            break


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


def calc_discount_price(product: NetseaProduct, discount_rate: float):
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
            netsea_object = NetseaProduct.get(url)
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


def netsea_all():
    new_shop_search()
    shops = NetseaShop.get_all_info()
    for shop in shops:
        run_netsea(shop.url)
