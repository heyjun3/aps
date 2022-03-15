import time
import queue
import threading
import json
from concurrent.futures import ThreadPoolExecutor

from spapi.spapi import SPAPI
from spapi.spapi import SPAPIJsonParser
from spapi.models import AsinsInfo
from spapi.models import SpapiFees
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
    with ThreadPoolExecutor(max_workers=6) as executor:
        for body in mq.get():
            executor.submit(threading_list_catalog_items, json.loads(body))

    logger.info('action=run_list_catalog_items status=done')

def threading_list_catalog_items(params: dict, interval_sec: float=0.17) -> None:
    logger.info('action=threading_list_catalog_items status=run')

    products = AsinsInfo.get(params['jan'])

    if not products:
        client = SPAPI()
        response = client.list_catalog_items(params['jan'])
        time.sleep(interval_sec)
        products = SPAPIJsonParser.parse_list_catalog_items(response.json())
        for product in products:
            AsinsInfo(asin=product['asin'], jan=params['jan'], title=product['title'], quantity=product['quantity']).upsert()
    else:
        products = [product.values for product in products]

    for product in products:
        MWS(asin=product['asin'], filename=params['filename'], title=product['title'],
                    jan=params['jan'], unit=product['quantity'], cost=params['cost']).save()
        
    logger.info('action=threading_list_catalog_items status=done')
    return None

def run_get_my_fees_estimate_for_asin() -> None:
    logger.info('action=run_get_my_fees_estimate_for_asin status=run')

    while True:
        asin_list = MWS.get_fee_is_None_asins()
        if asin_list:
            with ThreadPoolExecutor(max_workers=9) as executor:
                [executor.submit(threading_get_my_fees_estimate_for_asin, asin) for asin in asin_list]
        else:
            time.sleep(30)

def threading_get_my_fees_estimate_for_asin(asin: str, interval_sec: float=0.1, default_price: int=10000) -> None:
    logger.info('action=threading_get_my_fees_estimate_for_asin status=run')

    fee = SpapiFees.get(asin=asin)

    if fee is None:
        client = SPAPI()
        response = client.get_my_fees_estimate_for_asin(asin, price=default_price)
        time.sleep(interval_sec)
        fee = SPAPIJsonParser.parse_get_my_fees_estimate_for_asin(response.json())
        SpapiFees(asin=asin, fee_rate=fee['fee_rate'], shipping_fee=fee['ship_fee']).upsert()
    MWS.update_fee(asin=fee['asin'], fee_rate=fee['fee_rate'], shipping_fee=fee['ship_fee'])

    logger.info('action=threading_get_my_fees_estimate_for_asin status=done')
    return None


def main():
    client = SPAPI()
    asin_list = 'B000FQRA9S'
    response = client.get_item_offers(asin_list)
    # producSPAPIJsonParser.parse_get_competitive_pricing(response.json())
    print(response.json())