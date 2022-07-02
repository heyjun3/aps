import time
import queue
import threading
import datetime
from multiprocessing import Process
import json
from concurrent.futures import ThreadPoolExecutor
import asyncio

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

    def __init__(self, limit: int=20) -> None:
        self.queue = queue.Queue()
        self.asins = KeepaProducts.get_products_not_modified()
        self.asins = [self.asins[i:i+limit] for i in range(0, len(self.asins), limit)]
        self.spapi_client = SPAPI()

    async def main(self) -> None:
        get_competitive_pricing_task = asyncio.create_task(self.get_competitive_pricing())
        update_data_task = asyncio.create_task(self.update_data())

        asyncio.wait({get_competitive_pricing_task, update_data_task})

    async def update_data(self) -> None:

        while True:
            product = self.queue.get()
            if product is None:
                break

            now = time.time()
            KeepaProducts.update_price_and_rank_data(product['asin'], now, product['price'], product['ranking'])

    async def get_competitive_pricing(self, interval_sec: int=2) -> None:
        logger.info('action=get_competitive_pricing status=run')

        for asin_list in self.asins:
            response = await self.spapi_client.get_competitive_pricing(asin_list)
            products = SPAPIJsonParser.parse_get_competitive_pricing(response)
            [self.queue.put(product) for product in products]
            await asyncio.sleep(interval_sec)
        
        self.queue.put(None)
        logger.info('action=get_competitive_pricing status=done')


class RunAmzTask(object):

    def __init__(self, queue_name: str='mws', maxsize: int=10000) -> None:
        self.mq = MQ(queue_name)
        self.client = SPAPI()
        self.queue = queue.Queue(maxsize=maxsize)

    async def main(self) -> None:
        logger.info('action=main status=run')

        get_mq_task = asyncio.create_task(self.get_mq())
        search_catalog_items_task = asyncio.create_task(self.search_catalog_items_v20220401())
        await asyncio.wait({get_mq_task, search_catalog_items_task}, return_when='FIRST_COMPLETED')

        logger.info('action=main status=done')

    async def get_mq(self, interval_sec: float=2) -> None:
        logger.info('action=get_mq status=run')
        require = ('cost', 'jan', 'filename')
        mq_get_generator = self.mq.get()

        while True:
            params = await mq_get_generator.__anext__()
            if self.queue.full():
                await asyncio.sleep(interval_sec)
            params = json.loads(params)

            if not all(r in params for r in require):
                raise Exception

            products = AsinsInfo.get(params['jan'])
            products = False

            if not products:
                self.queue.put(params)
            else:
                for product in products:
                    MWS(asin=product['asin'], filename=params['filename'], title=product['title'],
                                jan=params['jan'], unit=product['quantity'], cost=params['cost']).save()

    async def search_catalog_items_v20220401(self, id_type: str='JAN', interval_sec: int=2) -> None:
        logger.info('action=search_catalog_items status=run')

        while True:
            params = [self.queue.get() for _ in range(20) if not self.queue.empty()]
            print({'queue_size': self.queue.qsize()})
            if not params:
                await asyncio.sleep(10)
            else:
                params = {param['jan']: {'filename': param['filename'], 'cost': param['cost']} for param in params}
                response = await self.client.search_catalog_items_v2022_04_01(params.keys(), id_type=id_type)
                products = SPAPIJsonParser.parse_search_catalog_items_v2022_04_01(response)
                for product in products:
                    parameter = params.get(product['jan'])
                    if parameter is None:
                        logger.error(product)
                        continue
                    AsinsInfo(asin=product['asin'], jan=product['jan'], title=product['title'], quantity=product['quantity']).upsert()
                    MWS(asin=product['asin'], filename=parameter['filename'], title=product['title'], jan=product['jan'], unit=product['quantity'], cost=parameter['cost']).save()

                await asyncio.sleep(interval_sec)

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
