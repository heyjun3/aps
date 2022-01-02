import time
from logging import getLogger
import configparser
import os
import pathlib
import shutil

import pandas as pd
import json
import requests

from models import KeepaProducts


logger = getLogger(__name__)
config = configparser.ConfigParser()
config.read(os.path.join(os.path.dirname(__file__), 'settings.ini'))
default = config['DEFAULT']

def check_keepa_tokens():
    logger.info('action=check_keepa_tokens status=run')
    url = f'https://api.keepa.com/product?key={default["KEEPA_ACCESS_KEY"]}'
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
        url = f'https://api.keepa.com/product?key={default["KEEPA_ACCESS_KEY"]}&domain=5&asin={asin_csv}&stats=90'

        while True:
            token = check_keepa_tokens()
            if token > len(asin_list):
                break
            time.sleep(60)

        response = requests(url)
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


def keepa_worker():
    logger.info('action=keepa_worker status=run')
    while True:
        try:
            path = next(pathlib.Path(config['AMAZON_RESULT_PATH']).iterdir())
        except StopIteration:
            logger.info('amazon_result_path is None')
            time.sleep(60)
            continue
        df = pd.read_pickle(str(path))
        drops = main(list(df['asin']))
        df = df.merge(drops, on='asin', how='inner').sort_values('drops', ascending=False).drop_duplicates()
        if not df.empty:
            df.to_excel(config['KEEPA_SAVE_PATH']+path.stem+'.xlsx', index=False)
        try:
            time.sleep(1)
            shutil.move(str(path), config['AMAZON_MOVE_PATH'])
        except Exception as e:
            logger.error(f'action=shutil.move error={e}')
            pass


if __name__ == '__main__':
    keepa_worker()
