import time
import queue
import threading

from spapi.spapi import SPAPI
from spapi.spapi import SPAPIJsonParser
from keepa.models import KeepaProducts
from mws.models import MWS
from mq import MQ
import log_settings


logger = log_settings.get_logger(__name__)


def update_price_and_ranking(interval_sec: int=2) -> None:
    logger.info('action=update_price_and_ranking status=run')

    que = queue.Queue()
    get_item_thread = threading.Thread(target=run_get_item_offers, args=(que, ))
    get_item_thread.start()

    products = KeepaProducts.get_products_not_modified()

    while products:
        asin_list = [products.pop() for _ in range(20) if products]
        thread = threading.Thread(target=run_get_competitive_pricing, args=(asin_list, que,))
        thread.start()
        time.sleep(interval_sec)

    que.put(None)
    get_item_thread.join()
    logger.info('action=update_price_and_ranking status=done')


def run_get_competitive_pricing(asin_list: list, que: queue.Queue) -> None:
    logger.info('action=run_parse_competitive_pricing status=run')

    client = SPAPI()
    response = client.get_competitive_pricing(asin_list)

    products = SPAPIJsonParser.parse_get_competitive_pricing(response.json())
    now = time.time()
    for product in products:
        KeepaProducts.update_price_and_rank_data(product['asin'], now, product['price'], product['ranking'])
        if product['price'] == -1:
            que.put(product['asin'])

    logger.info('action=run_parse_competitive_pricing status=done')


def run_get_item_offers(que: queue.Queue, interval_sec: float=0.2) -> None:
    logger.info('action=run_get_item_offers status=run')

    while True:
        asin = que.get()
        if asin is None:
            break
        thread = threading.Thread(target=get_item_offers_threading, args=(asin, ))
        thread.start()
        time.sleep(interval_sec)
       
    logger.info('action=run_get_item_offers status=done')


def get_item_offers_threading(asin: str) -> None:
    logger.info('action=get_item_offers_threading status=run')

    client = SPAPI()
    response = client.get_item_offers(asin)
    product = SPAPIJsonParser.parse_get_item_offers(response.json())
    if product is not None:
        KeepaProducts.update_price_and_rank_data(product['asin'], time.time(), product['price'], product['ranking'])

    logger.info('action=get_item_offers_threading status=done')


def run_list_catalog_items() -> None:
    logger.info('action=run_list_catalog_items status=run')

    mq = MQ('mws')
    mq.callback_recieve(func=threading_list_catalog_items, interval_sec=0.17)

    logger.info('action=run_list_catalog_items status=done')


def threading_list_catalog_items(params: dict) -> None:
    logger.info('action=threading_list_catalog_items status=run')

    client = SPAPI()
    response = client.list_catalog_items(params['jan'])
    products = SPAPIJsonParser.parse_list_catalog_items(response.json())
    if products:
        for product in products:
            mws = MWS(asin=product['asin'], filename=params['filename'], title=product['title'],
                      jan=params['jan'], unit=product['quantity'], price=product['price'], cost=params['cost'])
            mws.save()

    logger.info('action=threading_list_catalog_items status=done')


def main():
    client = SPAPI()
    asin_list = 'B000FQRA9S'
    response = client.get_item_offers(asin_list)
    # producSPAPIJsonParser.parse_get_competitive_pricing(response.json())
    print(response.json())