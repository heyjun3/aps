import time
import json
import itertools
import asyncio
import collections
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
        def _validation_parameter(param: dict) -> dict|None:
            if not all(r in param for r in require):
                logger.error({'bad parameter': param})
                return
            return param

        async def _get_asins_info_objects(jan_codes: List[str]) -> dict[str, List[AsinsInfo]]:
            asins = await AsinsInfo.get_asin_object_by_jan_list(jan_codes)
            asins = sorted(asins, key=lambda x: x.jan)
            return {k: list(g) for k, g in itertools.groupby(asins, lambda x: x.jan)}

        @log_decorator
        def _check_param_in_cache(param: dict) -> List[MWS]|None:
            if param is None: 
                return

            db_cache = asins.get(param['jan'])
            if db_cache is None:
                self.search_catalog_queue.publish(json.dumps(param))
                return
                
            return [MWS(
                asin=asin_info.asin,
                filename=param['filename'],
                title=asin_info.title,
                jan=asin_info.jan,
                unit=asin_info.quantity,
                cost=param['cost'],
                url=param['url'],
            ) for asin_info in db_cache]

        for messages in self.mq.receive(task_count):
            if messages is None:
                await asyncio.sleep(interval_sec)
                continue

            messages = [json.loads(message) for message in messages]
            jan_codes = list(map(lambda x: x['jan'], messages))
            asins = await _get_asins_info_objects(jan_codes)
            mws_objects = list(filter(None, reduce(lambda data, func: map(func, data),
                                [_validation_parameter, _check_param_in_cache], messages)))
            if mws_objects:
                mws_objects = itertools.chain.from_iterable(mws_objects)
                await MWS.insert_all_on_conflict_do_nothing(mws_objects)

    async def search_catalog_items_v20220401(self, id_type: str='JAN', interval_sec: int=2) -> None:
        logger.info('action=search_catalog_items status=run')

        for get_objects in self.search_catalog_queue.receive():
            if get_objects is None:
                logger.info({'message': 'get_objects is None'})
                await asyncio.sleep(10)
                continue

            params = sorted([json.loads(resp) for resp in get_objects], key=lambda x: x['jan'])
            params = {k: list(g) for k, g in itertools.groupby(params, lambda x: x['jan'])}

            response = await self.client.search_catalog_items_v2022_04_01(params.keys(), id_type=id_type)
            products = SPAPIJsonParser.parse_search_catalog_items_v2022_04_01(response)

            mws_objects = []
            for product in products:
                parameter = params.get(product['jan'])
                if parameter is None:
                    logger.error(product)
                    continue
                for param in parameter:
                    mws_objects.append(MWS(
                        asin=product['asin'],
                        filename=param['filename'],
                        title=product['title'],
                        jan=product['jan'],
                        unit=product['quantity'],
                        cost=param['cost'],
                        url=param['url'],
                    ))
            asyncio.ensure_future(AsinsInfo.insert_all_on_conflict_do_update(products))
            asyncio.ensure_future(MWS.insert_all_on_conflict_do_nothing(mws_objects))
            await asyncio.sleep(interval_sec)

    async def get_item_offers_batch(self, interval_sec: int=2):
        logger.info({'action': 'get_item_offers_batch', 'status': 'run'})
                       
        while True:
            mws_objects = await MWS.get_object_by_price_is_None()
            if not mws_objects:
                await asyncio.sleep(10)
                continue

            mws_objects = [mws_objects[i:i+20] for i in range(0, len(mws_objects), 20)]
            for mws_list in mws_objects:
                asins = [mws.asin for mws in mws_list]
                response = await self.client.get_item_offers_batch(asins)
                products = SPAPIJsonParser.parse_get_item_offers_batch(response)
                chain_products = collections.ChainMap(
                            *[{product['asin']: product['price']} for product in products])
                for mws in mws_list:
                    mws.price = chain_products.get(mws.asin, default=0)

                asyncio.ensure_future(MWS.insert_all_on_conflict_do_update_price(mws_list))
                asyncio.ensure_future(SpapiPrices.insert_all_on_conflict_do_update_price(products))
                await asyncio.sleep(interval_sec)

    async def get_my_fees_estimate(self, interval_sec: int=2) -> None:
        logger.info('action=get_my_fees_estimate_for_asin status=run')

        async def _get_my_fees_estimate(asin_list: List[str]) -> List[SpapiFees]:
            result = []
            if not asin_list:
                return

            for i in range(0, len(asin_list), 20):
                response = await self.client.get_my_fees_estimates(asin_list[i:i+20])
                products = SPAPIJsonParser.parse_get_my_fees_estimates(response)
                for product in products:
                    result.append(SpapiFees(product['asin'], product['fee_rate'], product['ship_fee']))
                await asyncio.sleep(interval_sec)
            return result

        while True:
            mws_objects = await MWS.get_fee_is_None_asins()
            if not mws_objects:
                await asyncio.sleep(30)
                continue

            asins = {mws.asin for mws in mws_objects}
            fees = await SpapiFees.get_asins_fee(list(asins))
            fees_asins = {fee.asin for fee in fees}
            search_asins = asins - fees_asins
            result = await _get_my_fees_estimate(list(search_asins))
            chain_map_fees = collections.ChainMap(
                                    *[{fee['asin']: fee} for fee in fees + result])

            for mws in mws_objects:
                fee = chain_map_fees.get(mws.asin)
                if fee:
                    mws.fee_rate = fee.fee_rate
                    mws.shipping_fee = fee.ship_fee
            asyncio.ensure_future(SpapiFees.insert_all_on_conflict_do_update_fee(result))
            asyncio.ensure_future(MWS.insert_all_on_conflict_do_update_fee(mws_objects))
