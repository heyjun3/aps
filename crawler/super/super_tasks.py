from datetime import datetime
from datetime import timedelta
import urllib.parse

import pandas as pd

from crawler.super.super import SuperCrawler
from crawler.super.models import SuperShop
import log_settings
import settings


logger = log_settings.get_logger(__name__)


def run_super_at_shop_id(shop_id: str):
    logger.info('action=run_super_at_shop_id status=run')

    url = urllib.parse.join(settings.SUPER_DOMAIN_URL, f'p/do/dpsl/{shop_id}')
    client = SuperCrawler(url=url)
    client.start_search_products()


def run_schedule_super_task():
    logger.info('action=run_schedule_super_task status=run')

    yesterday = datetime.now() - timedelta(days=1)
    url = urllib.parse.urljoin(settings.SUPER_NEW_PRODUCTS_URL, yesterday.strftime("%Y%m%d"))
    client = SuperCrawler(url=url)
    client.start_search_products

    logger.info('action=run_schedule_super_task status=done')


def run_discount_product_search():
    logger.info('action=run_discount_product_search status=run')

    url = urllib.parse.urljoin(settings.SUPER_DOMAIN_URL, 'p/do/psl/')
    params = {'pd': '1', 'is': '1', 'vi': '1'}
    client = SuperCrawler(url=url, params=params)
    client.start_search_products()

    logger.info('action=run_discount_product_search status=done')
    

def run_super_all_shops():
    logger.info('action=run_super_all_shops status=run')

    SuperShop.delete()
    run_get_super_shop_info()

    shops = SuperShop.get_all_info()
    for shop in shops:
        run_super_at_shop_id(shop.shop_id)

    logger.info('action=run_super_all_shops status=done')


def run_get_super_shop_info():
    logger.info('action=run_get_super_shop_info status=run')
    url = 'https://www.superdelivery.com/p/do/psl/?so=newdealer'

    client = SuperCrawler(url=url)
    client.pool_shop_list_page()

    logger.info('action=run_get_super_shop_info status=done')


def run_get_favorite_products() -> pd.DataFrame:
    logger.info('action=run_get_favorite_products status=run')
    
    url = 'https://www.superdelivery.com/p/wishlist/search.do'
    client = SuperCrawler(url=url)
    df = client.pool_favorite_product_list_page()

    logger.info('action=run_get_favorite_products status=done')
    return df