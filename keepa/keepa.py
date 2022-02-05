import time
from logging import getLogger
import os
import pathlib
import shutil
import datetime

import pandas as pd
import json
import requests

from keepa.models import KeepaProducts
import settings


logger = getLogger(__name__)


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


def check_keepa_tokens():
    logger.info('action=check_keepa_tokens status=run')
    url = f'https://api.keepa.com/token?key={settings.KEEPA_ACCESS_KEY}'
    response = request(url)
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

        token_count = check_keepa_tokens()
        if token_count < len(asin_list):
            interval_sec = (len(asin_list) - token_count) * 12 + 60
        else:
            interval_sec = 2
        time.sleep(interval_sec)
            

        response = request(url)
        response = response.json()

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
            
            KeepaProducts.update_or_insert(asin, drops, price_data, rank_data)
            data.append([asin, drops])

    df = pd.DataFrame(data=data, columns=['asin', 'drops']).astype({'drops': int})
    return df


def get_next_file_path():
    path = [path for path in pathlib.Path(settings.MWS_SAVE_PATH).iterdir()]
    if not path:
        return None
    else:
        path = sorted(path, key=lambda x: x.stat().st_mtime)
        return path[0]


def main(products: list):
    logger.info('action=main status=run')
    search_drop_list = []
    data = []

    for asin in products:
        db_object = KeepaProducts.object_get_db_asin(asin)

        if not db_object:
            search_drop_list.append(asin)
        elif db_object.price_data is None or db_object.rank_data is None:
            search_drop_list.append(asin)
        else:
            data.append([db_object.asin, db_object.sales_drops_90])

    result = keepa_get_drops(search_drop_list)
    df = pd.DataFrame(data=data, columns=['asin', 'drops']).astype({'drops': int})
    df = df.append(result, ignore_index=True)
    df = df.query('drops > 3')

    logger.info('action=main status=done')
    return df


def keepa_worker():
    logger.info('action=keepa_worker status=run')

    while True:
        path = get_next_file_path()
        if path is None:
            logger.info('amazon_result_path is None')
            tokens = check_keepa_tokens()
            if tokens > 100:
                products = KeepaProducts.get_product_price_data_is_None()
                if products:
                    asin_list = [product.asin for product in products]
                    keepa_get_drops(asin_list)
            time.sleep(60)
            continue
        else:
            df = pd.read_pickle(str(path))
            drops = main(list(df['asin']))
            df = df.merge(drops, on='asin', how='inner').sort_values('drops', ascending=False).drop_duplicates()
            if not df.empty:
                df.to_excel(os.path.join(settings.KEEPA_SAVE_PATH, f'{path.stem}.xlsx'), index=False)
            try:
                time.sleep(1)
                shutil.move(str(path), settings.MWS_DONE_SAVE_PATH)
            except Exception as e:
                logger.error(f'action=shutil.move error={e}')
                os.remove(str(path))
                pass
