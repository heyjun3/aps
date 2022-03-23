import time
import queue
import threading
from multiprocessing import Process
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
from mws import api


logger = log_settings.get_logger(__name__)


class UpdatePriceAndRankTask(object):

    def __init__(self, limit: int=20):
        self.queue = queue.Queue()
        self.asins = KeepaProducts.get_products_not_modified()
        self.asins = [self.asins[i:i+limit] for i in range(0, len(self.asins), limit)]
        self.spapi_client = SPAPI()

    def main(self):
        logger.info(f'action={self.__class__.__name__} main status=run')

        get_item_offers_thread = threading.Thread(target=self.get_item_offers_loop)
        get_competitive_pricing_thread = threading.Thread(target=self.get_competitive_pricing_loop)

        get_item_offers_thread.start()
        get_competitive_pricing_thread.start()

        get_competitive_pricing_thread.join()
        get_item_offers_thread.join()

        logger.info(f'action={self.__class__.__name__} main status=done')

    def get_competitive_pricing_loop(self, interval_sec: int=2):
        logger.info('action=get_competitive_pricing_loop status=run')

        with ThreadPoolExecutor(max_workers=10) as executor:
            for asin_list in self.asins:
                executor.submit(self.get_competitive_pricing, asin_list)
                time.sleep(interval_sec)

        self.queue.put(None)
        logger.info('action=get_competitive_pricing_loop status=done')

    def get_competitive_pricing(self, asin_list: list):
        logger.info('action=get_competitive_pricing status=run')

        response = self.spapi_client.get_competitive_pricing(asin_list)
        products = SPAPIJsonParser.parse_get_competitive_pricing(response.json())
        now = time.time()
        for product in products:
            KeepaProducts.update_price_and_rank_data(product['asin'], now, product['price'], product['ranking'])
            if product['price'] == -1:
                self.queue.put(product['asin'])

        logger.info('action=get_competitive_pricing status=done')

    def get_item_offers_loop(self, interval_sec: float=0.2):
        logger.info('action=get_item_offers_loop status=run')

        with ThreadPoolExecutor(max_workers=5) as executor:
            while True:
                asin = self.queue.get()
                if asin is None:
                    break
                executor.submit(self.get_item_offers, asin)
                time.sleep(interval_sec)

        logger.info('action=get_item_offers_loop status=done')

    def get_item_offers(self, asin: str):
        logger.info('action=get_item_offers status=run')

        response = self.spapi_client.get_item_offers(asin)
        product = SPAPIJsonParser.parse_get_item_offers(response.json())
        if product is not None:
            KeepaProducts.update_price_and_rank_data(product['asin'], time.time(), product['price'], product['ranking'])

        logger.info('action=get_item_offers status=done')

class RunAmzTask(object):

    def __init__(self, queue_name: str='mws'):
        self.mq = MQ(queue_name)
        self.client = SPAPI()

    def main(self):
        logger.info('action=main status=run')

        get_asins_info_process = Process(target=self.list_catalog_items_loop, daemon=True)
        get_price_process = Process(target=api.run_get_lowest_priced_offer_listtings_for_asin, daemon=True)
        get_fees_process = Process(target=self.get_my_fees_estimate_for_asin_loop, daemon=True)

        get_asins_info_process.start()
        get_price_process.start()
        get_fees_process.start()

        get_asins_info_process.join()
        get_price_process.join()
        get_fees_process.join()

        logger.info('action=main status=done')

    def list_catalog_items_loop(self) -> None:
        logger.info('action=list_catalog_items_loop status=run')

        with ThreadPoolExecutor(max_workers=6) as executor:
            for body in self.mq.get():
                executor.submit(self.list_catalog_items, json.loads(body))

        logger.info('action=list_catalog_items_loop status=done')


    def list_catalog_items(self, params: dict, interval_sec: float=0.17) -> None:
        logger.info('action=list_catalog_items status=run')

        products = AsinsInfo.get(params['jan'])

        if not products:
            response = self.client.list_catalog_items(params['jan'])
            time.sleep(interval_sec)
            products = SPAPIJsonParser.parse_list_catalog_items(response.json())
            for product in products:
                AsinsInfo(asin=product['asin'], jan=params['jan'], title=product['title'], quantity=product['quantity']).upsert()

        for product in products:
            MWS(asin=product['asin'], filename=params['filename'], title=product['title'],
                        jan=params['jan'], unit=product['quantity'], cost=params['cost']).save()
            
        logger.info('action=list_catalog_items status=done')
        return None


    def get_my_fees_estimate_for_asin_loop(self) -> None:
        logger.info('action=get_my_fees_estimate_for_asin_loop status=run')
        while True:
            asin_list = MWS.get_fee_is_None_asins()
            if asin_list:
                with ThreadPoolExecutor(max_workers=9) as executor:
                    [executor.submit(self.get_my_fees_estimate_for_asin, asin) for asin in asin_list]
            else:
                time.sleep(30)


    def get_my_fees_estimate_for_asin(self, asin: str, interval_sec: float=0.1, default_price: int=10000) -> None:
        logger.info('action=get_my_fees_estimate_for_asin status=run')

        fee = SpapiFees.get(asin=asin)

        if fee is None:
            response = self.client.get_my_fees_estimate_for_asin(asin, price=default_price)
            time.sleep(interval_sec)
            fee = SPAPIJsonParser.parse_get_my_fees_estimate_for_asin(response.json())
            SpapiFees(asin=asin, fee_rate=fee['fee_rate'], ship_fee=fee['ship_fee']).upsert()
        MWS.update_fee(asin=fee['asin'], fee_rate=fee['fee_rate'], shipping_fee=fee['ship_fee'])

        logger.info('action=get_my_fees_estimate_for_asin status=done')
        return None


def main():
    client = SPAPI()
    asin_list = 'B000FQRA9S'
    response = client.get_item_offers(asin_list)
    # producSPAPIJsonParser.parse_get_competitive_pricing(response.json())
    print(response.json())