import time

import json
import requests

from keepa.models import KeepaProducts
from mws.models import MWS
import settings
import log_settings


logger = log_settings.get_logger(__name__)


def request(url):
    for _ in range(60):
        try:
            response = requests.post(url)
            if response.status_code == 200:
                return response
            else:
                logger.error(f'Request Error code {response.status_code}')
                raise Exception
        except Exception as ex:
            logger.error(ex)
            time.sleep(60)


def check_keepa_tokens(interval_sec: int = 1):
    logger.info('action=check_keepa_tokens status=run')
    url = f'https://api.keepa.com/token?key={settings.KEEPA_ACCESS_KEY}'
    response = request(url)
    time.sleep(interval_sec)
    response = response.content.decode()
    response = json.loads(response)
    logger.info(f'tokens:{response["tokensLeft"]}')
    return int(response["tokensLeft"])


# argument: asin list ,Return: [asin, 90drops]
def keepa_request_products(asin_list: list) -> dict:
    logger.info('action=keepa_get_drops status=run')

    asin_csv = ','.join(asin_list)
    url = f'https://api.keepa.com/product?key={settings.KEEPA_ACCESS_KEY}&domain=5&asin={asin_csv}&stats=90'

    token_count = check_keepa_tokens()
    if token_count < len(asin_list):
        interval_sec = (len(asin_list) - token_count) * 12 + 60
    else:
        interval_sec = 2
    
    logger.info('Waiting recovery keepa tokens')
    time.sleep(interval_sec)
        
    response = request(url)
    response = response.json()

    return response


def scrape_keepa_request(response: dict) -> list:
    logger.info('action=scrape_keepa_request status=run')
    products = []
    PRICE_DATA_NUM = 1
    RANK_DATA_NUM = 3

    for product in response.get('products'):
        asin = product.get('asin')
        drops = int(product.get('stats').get('salesRankDrops90'))

        try:
            price_data = product.get('csv')[PRICE_DATA_NUM]
            price_data = {date: price for date, price in zip(price_data[0::2], price_data[1::2])}
        except TypeError as ex:
            logger.error(f"{asin} hasn't price data {ex}")
            price_data = {'-1': -1}
        try:
            rank_data = product.get('csv')[RANK_DATA_NUM]
            rank_data = {date: rank for date, rank in zip(rank_data[0::2], rank_data[1::2])}
        except TypeError as ex:
            logger.error(f"{asin} hasn't rank data {ex}")
            rank_data = {'-1': -1}
        products.append((asin, drops, price_data, rank_data))

    return products


async def main(interval_sec: int=60, count: int=100):
    logger.info('action=main status=run')

    while True:
        asin_list = await MWS.get_asin_list_None_products()
        if asin_list:
            asin_list = [asin_list[i:i+count] for i in range(0, len(asin_list), count)]
            for asins in asin_list:
                response = keepa_request_products(asins)
                keepa_products = scrape_keepa_request(response)
                for asin, drops, price_data, rank_data in keepa_products:
                    KeepaProducts.update_or_insert(asin, drops, price_data, rank_data)
        else:
            time.sleep(interval_sec)
