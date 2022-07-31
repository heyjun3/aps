import time
from multiprocessing import Process, Queue
import json
import asyncio
from typing import Coroutine

from spapi.spapi import SPAPI
from spapi.spapi import SPAPIJsonParser
from spapi.models import AsinsInfo, SpapiPrices
from spapi.models import SpapiFees
from keepa.models import KeepaProducts
from mws.models import MWS
from mq import MQ
import log_settings


logger = log_settings.get_logger(__name__)


class UpdatePriceAndRankTask(object):

    def __init__(self) -> None:
        self.queue = Queue()
        self.spapi_client = SPAPI()

    async def main(self, limit: int=20, timeout: int=86400) -> None:
        logger.info('action=main status=run')
        update_data_process = Process(target=self.update_data, args=(self.queue, ))
        update_data_process.start()

        try:
            while True:
                self.asins = KeepaProducts.get_products_not_modified()
                if not self.asins:
                    self.queue.put(None)
                    break
                self.asins = [self.asins[i:i+limit] for i in range(0, len(self.asins), limit)]
                get_competitive_pricing_task = asyncio.create_task(self.get_competitive_pricing())
                await asyncio.wait_for(get_competitive_pricing_task, timeout=timeout)
            update_data_process.join()
        except asyncio.TimeoutError as ex:
            logger.error(f'action=main error={ex}')
            self.queue.put(None)
            update_data_process.join()

        logger.info('action=main status=done')

    def update_data(self, queue: Queue) -> None:
        logger.info('action=update_data status=run')

        while True:
            product = queue.get()
            if product is None:
                break

            now = time.time()
            KeepaProducts.update_price_and_rank_data(product['asin'], now, product['price'], product['ranking'])

        logger.info('action=update_data status=done')

    async def get_competitive_pricing(self, interval_sec: int=2) -> None:
        logger.info('action=get_competitive_pricing status=run')

        async def _get_competitive_pricing(asin_list):
            response = await self.spapi_client.get_competitive_pricing(asin_list)
            products = SPAPIJsonParser.parse_get_competitive_pricing(response)
            [self.queue.put(product) for product in products]

        for asin_list in self.asins:
            task = asyncio.create_task(_get_competitive_pricing(asin_list))
            sleep = asyncio.create_task(asyncio.sleep(interval_sec))
            await asyncio.gather(task, sleep)
        
        logger.info('action=get_competitive_pricing status=done')


class RunAmzTask(object):

    def __init__(self, queue_name: str='mws', search_queue: str='search_catalog') -> None:
        self.mq = MQ(queue_name)
        self.search_catalog_queue = MQ(search_queue)
        self.client = SPAPI()
        self.fees_queue = asyncio.Queue()

    async def main(self) -> None:
        logger.info('action=main status=run')

        get_mq_process = Process(target=asyncio.run, args=(self.get_queue(), ))
        search_catalog_items_process = Process(target=asyncio.run, args=(self.search_catalog_items_v20220401(), ))
        get_item_offers_process = Process(target=asyncio.run, args=(self.get_item_offers_batch(), ))
        get_my_fee_estimate_process = Process(target=asyncio.run, args=(self.get_my_fees_estimate(), ))

        get_mq_process.start()
        search_catalog_items_process.start()
        get_item_offers_process.start()
        get_my_fee_estimate_process.start()

        get_mq_process.join()
        search_catalog_items_process.join()
        get_item_offers_process.join()
        get_my_fee_estimate_process.join()

        logger.info('action=main status=done')

    async def get_queue(self, interval_sec: int=10) -> None:
        logger.info('action=get_queue status=run')

        require = ('cost', 'jan', 'filename')

        for params in self.mq.get():
            if params is None:
                await asyncio.sleep(interval_sec)
                continue

            params_json = json.loads(params)
            if not all(r in params_json for r in require):
                logger.error({'bad parameter': params_json})
                continue

            products = await AsinsInfo.get(params_json['jan'])

            if not products:
                self.search_catalog_queue.publish(params)
            else:
                for product in products:
                    await MWS(asin=product['asin'], filename=params_json['filename'], title=product['title'],
                        jan=params_json['jan'], unit=product['quantity'], cost=params_json['cost']).save()

    async def search_catalog_items_v20220401(self, id_type: str='JAN', interval_sec: int=2) -> None:
        logger.info('action=search_catalog_items status=run')

        while True:
            params = [next(self.search_catalog_queue.get()) for _ in range(20)]
            params = [json.loads(param) for param in params if param is not None]

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
                    await AsinsInfo(asin=product['asin'], jan=product['jan'], title=product['title'], quantity=product['quantity']).upsert()
                    await MWS(asin=product['asin'], filename=parameter['filename'], title=product['title'], jan=product['jan'], unit=product['quantity'], cost=parameter['cost']).save()

                await asyncio.sleep(interval_sec)

    async def get_item_offers_batch(self, interval_sec: int=2):
        logger.info({'action': 'get_item_offers_batch', 'status': 'run'})

        async def _get_item_offers_batch(asins):
            logger.info({'action': '_get_item_offers_batch', 'status': 'run'})
            response = await self.client.get_item_offers_batch(asins)
            products = SPAPIJsonParser.parse_get_item_offers_batch(response)
            for product in products:
                await MWS.update_price(asin=product['asin'], price=product['price'])
                await SpapiPrices(asin=product['asin'], price=product['price']).upsert()
            
            logger.info({'action': '_get_item_offers_batch', 'status': 'done'})

        while True:
            asin_list = await MWS.get_price_is_None_asins()
            if asin_list:
                asin_list = [asin_list[i:i+20] for i in range(0, len(asin_list), 20)]
                for asins in asin_list:
                    task = asyncio.create_task(_get_item_offers_batch(asins))
                    sleep = asyncio.create_task(asyncio.sleep(interval_sec))
                    await asyncio.gather(task, sleep)
            else:
                await asyncio.sleep(10)

    async def get_my_fees_estimate(self, interval_sec: int=2, use_cache: bool=True) -> None:
        logger.info('action=get_my_fees_estimate_for_asin status=run')

        async def _get_my_fees_estimate(asin_list):
            async def _wrapper_get_my_fees_estimate(asins):
                response = await self.client.get_my_fees_estimates(asins)
                products = SPAPIJsonParser.parse_get_my_fees_estimates(response)
                for product in products:
                    await SpapiFees(asin=product['asin'], fee_rate=product['fee_rate'], ship_fee=product['ship_fee']).upsert()
                    await MWS.update_fee(asin=product['asin'], fee_rate=product['fee_rate'], shipping_fee=product['ship_fee'])

            asin_list = [asin_list[i:i+20] for i in range(0, len(asin_list), 20)]
            for asins in asin_list:
                task = asyncio.create_task(_wrapper_get_my_fees_estimate(asins))
                sleep = asyncio.create_task(asyncio.sleep(interval_sec))
                await asyncio.gather(task, sleep)

        while True:
            asin_list = await MWS.get_fee_is_None_asins()
            if asin_list:
                fees = []
                for asin in asin_list:
                    fee = await SpapiFees.get(asin)
                    if fee is None:
                        fees.append(asin)
                    else:
                        await MWS.update_fee(fee['asin'], fee['fee_rate'], fee['ship_fee'])

                await _get_my_fees_estimate(fees)    
            else:
                await asyncio.sleep(30)
