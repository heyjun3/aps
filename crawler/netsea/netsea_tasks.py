import urllib.parse
import time
import datetime

import pandas as pd
import requests

from crawler.netsea.netsea import Netsea
from crawler.netsea.netsea import NetseaHTMLPage
from crawler.netsea.models import NetseaShop
from crawler import utils
import settings
import log_settings


logger = log_settings.get_logger(__name__)

def logger_decorator(func):
    def _logger_decorator(*args, **kwargs):
        logger.info({'action': func.__name__, 'status': 'run'})
        result = func(*args, **kwargs)
        logger.info({'action': func.__name__, 'status': 'done'})
        return result
    return _logger_decorator

@logger_decorator
def run_netsea_at_shop_id(shop_id: str, path: str = 'search_faceted') -> None:
    url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)
    params = {'sort': 'PD', 'supplier_id': shop_id, 'ex_so': 'Y', 'searched': 'Y'}
    url = requests.Request(method='GET', url=url, params=params).prepare().url
    timestamp = datetime.datetime.now()
    client = Netsea([url], timestamp=timestamp)
    client.start_search_products()


@logger_decorator
def run_new_product_search(path: str = 'search_faceted') -> None:
    timestamp = datetime.datetime.now()
    base_url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)
    urls = []

    for index in range(2, 8):
        params = {'sort': 'new', 'category_id': str(index), 'ex_so': 'Y', 'searched': 'Y'}
        url = requests.Request(method='GET', url=base_url, params=params).prepare().url
        urls.append(url)
    
    client = Netsea(urls, timestamp, is_new_product_search=True)
    client.start_search_products()

@logger_decorator
def run_get_discount_products(path: str = 'search_faceted') -> None:
    base_url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)
    timestamp = datetime.datetime.now()
    urls = []

    for index in range(10, 1, -1):
        params = {'disc_flg': 'Y', 'ex_so': 'Y', 'sort': 'PD', 'searched': 'Y', 'category_id': str(index)}
        url = requests.Request(method='GET', url=base_url, params=params).prepare().url
        urls.append(url)

    client = Netsea(urls, timestamp=timestamp)
    client.start_search_products()

@logger_decorator
def run_netsea_all_products():
    NetseaShop.delete()
    run_get_all_shop_info()
    shops = NetseaShop.get_all_info()
    for shop in shops:
        logger.info({'shop_id': shop.shop_id})
        run_netsea_at_shop_id(shop.shop_id)

@logger_decorator
def run_get_all_shop_info(path: str = 'shop', interval_sec: int = 2) -> None:
    url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)

    for index in range(1, 9):
        params = {'category_id': str(index), 'sort': 'NEW'}
        response = utils.request(url=url, params=params)
        shops = NetseaHTMLPage.scrape_shop_list_page(response.text)
        list(map(lambda x: x.save(), shops))
        time.sleep(interval_sec)

@logger_decorator
def run_get_favorite_products(path: str = 'bookmark') -> pd.DataFrame:
    url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)
    params = {'stock_option': 'in'}
    url = requests.Request(method='GET', url=url, params=params).prepare().url
    client = Netsea([url], params)
    df = client.pool_favorite_product_list_page()
    return df
