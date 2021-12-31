import time
import datetime
from logging import getLogger
import logging.handlers
from collections import deque
import logging.config

import pandas as pd
import json
import requests

from scraping.controllers.utils import post_request
from scraping.models.models import KeepaProducts
import settings

logger = getLogger(__name__)


def check_keepa_tokens():
    logger.info('action=check_keepa_tokens status=run')
    url = f'https://api.keepa.com/product?key={settings.KEEPA_ACCESS_KEY}'
    response = requests.post(url)
    time.sleep(1)
    response = response.content.decode()
    response = json.loads(response)
    logger.info(f'tokens:{response["tokensLeft"]}')
    return int(response["tokensLeft"])


# argument: asin list ,Return: [asin, 90drops]
def keepa_get_drops(products: list):
    logger.info('action=keepa_get_drops status=run')
    data = []

    while products:
        asin_list = [products.pop() for _ in range(100) if products]
        asin_csv = ','.join(asin_list)
        url = f'https://api.keepa.com/product?key={settings.KEEPA_ACCESS_KEY}&domain=5&asin={asin_csv}&stats=90'

        while True:
            token = check_keepa_tokens()
            if token > len(asin_list):
                break
            time.sleep(60)

        response = post_request(url)
        response = response.json()
        time.sleep(2)

        for product in response.get('products'):
            asin = product.get('asin')
            drops = int(product.get('stats').get('salesRankDrops90'))
            KeepaProducts.update_or_insert(asin, drops)
            data.append([asin, drops])

    df = pd.DataFrame(data=data, columns=['asin', 'drops']).astype({'drops': int})
    return df


def main(products: list):
    logger.info('action=main status=run')
    search_drop_list = []
    data = []

    for asin in products:
        db_object = KeepaProducts.object_get_db_asin(asin)

        if not db_object:
            search_drop_list.append(asin)
        else:
            data.append([db_object.asin, db_object.sales_drops_90])

    result = keepa_get_drops(search_drop_list)
    df = pd.DataFrame(data=data, columns=['asin', 'drops']).astype({'drops': int})
    df = df.append(result, ignore_index=True)
    df = df.query('drops > 3')

    logger.info('action=main status=done')
    return df


# if __name__ == '__main__':

