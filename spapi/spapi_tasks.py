import time
import queue
import threading

from spapi.spapi import SPAPI
from spapi.spapi import SPAPIJsonParser
from keepa.models import KeepaProducts
import log_settings


logger = log_settings.get_logger(__name__)


def get_landed_price_and_ranking(interval_sec: int=2) -> None:
    logger.info('action=get_landed_price_and_ranking status=run')

    que = queue.Queue()
    thread = threading.Thread(target=run_parse_competitive_pricing, args=(que, ))
    thread.start()
    client = SPAPI()

    while True:
        products = KeepaProducts.get_products_not_modified()
        if not products:
            break

        asin_list = [product[0] for product in products]
        response = client.get_competitive_pricing(asin_list)
        que.put(response)
        time.sleep(interval_sec)


    que.put(None)
    thread.join()
    logger.info('action=get_landed_price_and_ranking status=done')


def run_parse_competitive_pricing(get_que: queue.Queue) -> None:
    logger.info('action=run_parse_competitive_pricing status=run')

    put_que = queue.Queue()
    thread = threading.Thread(target=run_get_item_offers, args=(put_que, ))
    thread.start()

    while True:
        response = get_que.get()

        if response is None:
            break
        else:
            products = SPAPIJsonParser.parse_get_competitive_pricing(response.json())
            now = time.time()
            for product in products:
                KeepaProducts.update_price_and_rank_data(product['asin'], now, product['price'], product['ranking'])
                if product['price'] == -1:
                    put_que.put(product['asin'])
    put_que.put(None)
    logger.info('action=run_parse_competitive_pricing status=done')


def run_get_item_offers(que: queue.Queue, interval_sec: float=0.2) -> None:
    logger.info('action=run_get_item_offers status=run')

    client = SPAPI()

    while True:
        asin = que.get()
        if asin is None:
            break
        response = client.get_item_offers(asin)
        time.sleep(interval_sec)
        product = SPAPIJsonParser.parse_get_item_offers(response.json())
        if product is not None:
            KeepaProducts.update_price_and_rank_data(product['asin'], time.time(), product['price'], product['ranking'])

    logger.info('action=run_get_item_offers status=done')


def main():
    client = SPAPI()
    asin_list = 'B000FQRA9S'
    response = client.get_item_offers(asin_list)
    # producSPAPIJsonParser.parse_get_competitive_pricing(response.json())
    print(response.json())