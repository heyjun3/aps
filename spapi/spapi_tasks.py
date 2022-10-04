import time
import json
import itertools
import asyncio
from typing import List
from typing import Callable
from multiprocessing import Process, Queue
from functools import reduce
from functools import partial

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


def log_decorator(func: Callable) -> Callable:
    def _inner(*args, **kwargs) -> any:
        logger.info({'action': func.__name__, 'status': 'run'})
        result = func(*args, **kwargs)
        logger.info({'action': func.__name__, 'status': 'done'})
        return result
    return _inner


class UpdatePriceAndRankTask(object):

    def __init__(self) -> None:
        self.queue = Queue()
        self.spapi_client = SPAPI()

    async def main(self, limit: int=20) -> None:
        logger.info('action=main status=run')
        update_data_process = Process(target=self.update_data, args=(self.queue, ))
        update_data_process.start()

        while True:
            self.asins = KeepaProducts.get_products_not_modified()
            if not self.asins:
                time.sleep(60)
                continue
            self.asins = [self.asins[i:i+limit] for i in range(0, len(self.asins), limit)]
            get_competitive_pricing_task = asyncio.create_task(self.get_competitive_pricing())
            await get_competitive_pricing_task

    def update_data(self, queue: Queue) -> None:
        logger.info('action=update_data status=run')

        while True:
            product = queue.get()
            now = time.time()
            KeepaProducts.update_price_and_rank_data(product['asin'], now, product['price'], product['ranking'])

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
        self.jan_cache = Cache(None, 3600)
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

        require = ('cost', 'jan', 'filename', 'url')

        @log_decorator
        def _validation_parameter(parameter: str) -> dict|None:
            param = json.loads(parameter)
            if not all(r in param for r in require):
                logger.error({'bad parameter': param})
                return
            return param

        @log_decorator
        def _check_param_in_cache(param: dict) -> dict|None:
            if param is None: 
                return

            if param['jan'] in jan_cache:
                return param

            self.search_catalog_queue.publish(json.dumps(param))

        @log_decorator
        def _combine_param_and_asins_objects(asin_object: AsinsInfo, params: List[dict]) -> MWS:
            filter_params = filter(lambda x: x['jan'] == asin_object.jan, params)
            return [MWS(
                asin=asin_object.asin,
                filename = param['filename'],
                title=asin_object.title,
                jan=asin_object.jan,
                unit=asin_object.quantity,
                cost=param['cost'],
                url=param['url']) for param in filter_params]

        for messages in self.mq.receive(task_count):
            if messages is None:
                await asyncio.sleep(interval_sec)
                continue

            jan_cache = self.jan_cache.get_value()
            if jan_cache is None:
                jan_cache = await AsinsInfo.get_jan_code_all()
                self.jan_cache.set_value(set(jan_cache))

            messages = list(filter(None, reduce(lambda data, func: map(func, data), [_validation_parameter, _check_param_in_cache], messages)))
            if messages:
                jan_list = list(map(lambda x: x.get('jan'), messages))
                asins = await AsinsInfo.get_asin_object_by_jan_list(jan_list)
                mws_records = map(partial(_combine_param_and_asins_objects, params=messages), asins)
                mws_records = list(itertools.chain.from_iterable(mws_records))
                await MWS.insert_all_on_conflict_do_nothing(mws_records)

    async def search_catalog_items_v20220401(self, id_type: str='JAN', interval_sec: int=2) -> None:
        logger.info('action=search_catalog_items status=run')

        for get_objects in self.search_catalog_queue.receive():
            if get_objects is None:
                logger.info({'message': 'get_objects is None'})
                await asyncio.sleep(10)
                continue

            params = [json.loads(resp) for resp in get_objects]

            params = {param['jan']: {'filename': param['filename'], 'cost': param['cost'], 'url': param['url']} for param in params}
            response = await self.client.search_catalog_items_v2022_04_01(params.keys(), id_type=id_type)
            products = SPAPIJsonParser.parse_search_catalog_items_v2022_04_01(response)
            for product in products:
                parameter = params.get(product['jan'])
                if parameter is None:
                    logger.error(product)
                    continue
                asyncio.ensure_future(AsinsInfo(asin=product['asin'], jan=product['jan'], title=product['title'], quantity=product['quantity']).upsert())
                asyncio.ensure_future(MWS(asin=product['asin'], filename=parameter['filename'], title=product['title'], jan=product['jan'], unit=product['quantity'], cost=parameter['cost'], url=parameter['url']).save())

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

            fees = await SpapiFees.get_asins_fee(asins)
            for fee in fees:
                await MWS.update_fee(fee['asin'], fee['fee_rate'], fee['ship_fee'])

        async def _get_my_fees_estimate(asin_list: List[str]) -> None:
            if not asin_list:
                return

            asin_collection = [asin_list[i:i+20] for i in range(0, len(asin_list), 20)]
            for asins in asin_collection:
                response = await self.client.get_my_fees_estimates(asins)
                products = SPAPIJsonParser.parse_get_my_fees_estimates(response)
                for product in products:
                    asyncio.ensure_future(SpapiFees(asin=product['asin'], fee_rate=product['fee_rate'], ship_fee=product['ship_fee']).upsert())
                    asyncio.ensure_future(MWS.update_fee(asin=product['asin'], fee_rate=product['fee_rate'], shipping_fee=product['ship_fee']))
                await asyncio.sleep(interval_sec)

        while True:
            asin_list = await MWS.get_fee_is_None_asins()
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

            insert_task = asyncio.create_task(_insert_db_using_cache(asins_exist_db))
            get_info_task = asyncio.create_task(_get_my_fees_estimate(asin_list))
            await asyncio.gather(insert_task, get_info_task)
