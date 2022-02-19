import urllib.parse
import time
import datetime

import pandas as pd

from crawler.netsea.netsea import Netsea
from crawler.netsea.netsea import NetseaHTMLPage
from crawler.netsea.models import NetseaShop
from crawler import utils
import settings
import log_settings


logger = log_settings.get_logger(__name__)


def run_netsea_at_shop_id(shop_id: str, path: str = 'search') -> None:
    url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)
    params = {'sort': 'PD', 'supplier_id': shop_id, 'ex_so': 'Y', 'searched': 'Y'}
    client = Netsea(url, params)
    client.start_search_products()


def run_new_product_search(path: str = 'search') -> None:
    logger.info('action=run_new_product_search status=run')

    timestamp = datetime.datetime.now()
    url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)

    for index in range(1, 10):
        params = {'sort': 'new', 'category_id': str(index), 'ex_so': 'Y'}
        client = Netsea(url, params, timestamp)
        client.start_search_products()


def run_get_discount_products(path: str = 'search') -> None:
    logger.info('action=run_get_discount_products')

    url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)
    params = {'sort': 'PD', 'ex_so': 'Y', 'searhed': 'Y', 'disc_flg': 'Y'}
    client = Netsea(url, params)
    client.start_search_products()


def run_netsea_all_products():
    NetseaShop.delete()
    run_get_all_shop_info()
    shops = NetseaShop.get_all_info()
    for shop in shops:
        run_netsea_at_shop_id(shop.shop_id)


def run_get_all_shop_info(path: str = 'shop', interval_sec: int = 2) -> None:
    logger.info('action=run_get_all_shop_info status=run')
    url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)

    for index in range(1, 9):
        params = {'category_id': str(index), 'sort': 'NEW'}
        response = utils.request(url=url, params=params)
        shops = NetseaHTMLPage.scrape_shop_list_page(response.text)
        list(map(lambda x: x.save(), shops))
        time.sleep(interval_sec)


def run_get_favorite_products(path: str = 'bookmark') -> pd.DataFrame:

    url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)
    params = {'stock_option': 'in'}
    client = Netsea(url, params)
    df = client.pool_favorite_product_list_page()
    return df
