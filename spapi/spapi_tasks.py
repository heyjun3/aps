import time

from spapi.spapi import SPAPI
from spapi.spapi import SPAPIJsonParser
from keepa.models import KeepaProducts
import log_settings


logger = log_settings.get_logger(__name__)


def get_landed_price_and_ranking(interval_sec: int=2) -> None:
    logger.info('action=get_landed_price_and_ranking status=run')

    client = SPAPI()

    while True:
        products = KeepaProducts.get_products_not_modified()
        if not products:
            break

        asin_list = [product.asin for product in products]
        response = client.get_competitive_pricing(asin_list)
        time.sleep(interval_sec)
        products = SPAPIJsonParser.parse_get_competitive_pricing(response.json())
        now = time.time()
        for product in products:
            KeepaProducts.update_price_and_rank_data(product['asin'], now, product['price'], product['ranking'])
    
    logger.info('action=get_landed_price_and_ranking status=done')


def main():
    client = SPAPI()
    asin_list = ['B08HMT3LRN', 'B07HG6F6K2']
    response = client.get_competitive_pricing(asin_list)
    # producSPAPIJsonParser.parse_get_competitive_pricing(response.json())
    print(response.json())