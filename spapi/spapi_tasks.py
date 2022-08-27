import time
from multiprocessing import Process, Queue
import json
import asyncio
import functools
from typing import List

from spapi.spapi import SPAPI
from spapi.spapi import SPAPIJsonParser
from spapi.models import AsinsInfo, SpapiPrices
from spapi.models import SpapiFees
from keepa.models import KeepaProducts
from mws.models import MWS
from mq import MQ
from spapi.utils import Cache
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
        self.cache = Cache(None, 3600)
        self.estimate_queue = asyncio.Queue()

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

    async def get_queue(self, interval_sec: int=10, task_count=100) -> None:
        logger.info('action=get_queue status=run')

        require = ('cost', 'jan', 'filename')

        async def _get_queue():
            params = self.mq.basic_get()
            if params is None:
                await asyncio.sleep(interval_sec)
                return

            params_json = json.loads(params)
            if not all(r in params_json for r in require):
                logger.error({'bad parameter': params_json})
                return

            products = await AsinsInfo.get(params_json['jan'])
            if not products:
                self.search_catalog_queue.publish(params)
            else:
                for product in products:
                    asyncio.ensure_future(MWS(asin=product['asin'], filename=params_json['filename'], title=product['title'],
                        jan=params_json['jan'], unit=product['quantity'], cost=params_json['cost']).save())

        while True:
            if self.mq.get_message_count:
                tasks = [asyncio.create_task(_get_queue()) for _ in range(task_count)]
                await asyncio.gather(*tasks)
            else:
                asyncio.sleep(interval_sec)

    async def search_catalog_items_v20220401(self, id_type: str='JAN', interval_sec: int=2) -> None:
        logger.info('action=search_catalog_items status=run')

        for get_objects in self.search_catalog_queue.receive():
            if get_objects is None:
                logger.info({'message': 'get_objects is None'})
                await asyncio.sleep(10)
                continue

            params = [json.loads(resp) for resp in get_objects]

            params = {param['jan']: {'filename': param['filename'], 'cost': param['cost']} for param in params}
            response = await self.client.search_catalog_items_v2022_04_01(params.keys(), id_type=id_type)
            products = SPAPIJsonParser.parse_search_catalog_items_v2022_04_01(response)
            for product in products:
                parameter = params.get(product['jan'])
                if parameter is None:
                    logger.error(product)
                    continue
                asyncio.ensure_future(AsinsInfo(asin=product['asin'], jan=product['jan'], title=product['title'], quantity=product['quantity']).upsert())
                asyncio.ensure_future(MWS(asin=product['asin'], filename=parameter['filename'], title=product['title'], jan=product['jan'], unit=product['quantity'], cost=parameter['cost']).save())

            await asyncio.sleep(interval_sec)

    async def get_item_offers_batch(self, interval_sec: int=2):
        logger.info({'action': 'get_item_offers_batch', 'status': 'run'})
                       
        while True:
            asin_list = await MWS.get_price_is_None_asins()
            if asin_list:
                asin_list = [asin_list[i:i+20] for i in range(0, len(asin_list), 20)]
                for asins in asin_list:
                    response = await self.client.get_item_offers_batch(asins)
                    products = SPAPIJsonParser.parse_get_item_offers_batch(response)
                    for product in products:
                        asyncio.ensure_future(MWS.update_price(asin=product['asin'], price=product['price']))
                        asyncio.ensure_future(SpapiPrices(asin=product['asin'], price=product['price']).upsert())
                    
                    await asyncio.sleep(interval_sec)
            else:
                await asyncio.sleep(10)

    async def get_my_fees_estimate(self, interval_sec: int=2) -> None:
        logger.info('action=get_my_fees_estimate_for_asin status=run')

        async def _insert_db_using_cache(asins: List[str]) -> None:

            for asin in asins:
                fee = await SpapiFees.get(asin)
                logger.info(fee)
                await MWS.update_fee(fee['asin'], fee['fee_rate'], fee['ship_fee'])

        async def _get_my_fees_estimate(asin_list: List[str]) -> None:
            if not asin_list:
                return

            asin_collection = [asin_list[i:i+20] for i in range(0, len(asin_list), 20)]
            for asins in asin_collection:
                response = await self.client.get_my_fees_estimates(asins)
                products = SPAPIJsonParser.parse_get_my_fees_estimates(response)
                logger.info(products)
                for product in products:
                    asyncio.ensure_future(SpapiFees(asin=product['asin'], fee_rate=product['fee_rate'], ship_fee=product['ship_fee']).upsert())
                    asyncio.ensure_future(MWS.update_fee(asin=product['asin'], fee_rate=product['fee_rate'], shipping_fee=product['ship_fee']))
                await asyncio.sleep(interval_sec)

        while True:
            asin_list = await MWS.get_fee_is_None_asins(10)
            if not asin_list:
                await asyncio.sleep(30)
                continue

            if self.cache.get_value() is None:
                asins = await SpapiFees.get_asins_after_update_interval_days()
                self.cache.set_value(set(asins))

            asins_in_database = self.cache.get_value()

            asins_exist_db = []
            asin_list = [asins_exist_db.append(asin) if asin in asins_in_database else asin for asin in asin_list]
            asin_list = [asin for asin in asin_list if asin]
            logger.info(asin_list)
            logger.info(asins_exist_db)

            insert_task = asyncio.create_task(_insert_db_using_cache(asins_exist_db))
            get_info_task = asyncio.create_task(_get_my_fees_estimate(asin_list))
            await asyncio.gather(insert_task, get_info_task)
            break
