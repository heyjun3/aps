from copy import deepcopy
import multiprocessing
import time
import json
import itertools
import asyncio
import collections
import re
from typing import ChainMap, List
from typing import Callable
from multiprocessing import Process, Queue
from functools import reduce
from functools import partial

from spapi.spapi import SPAPI
from spapi.spapi import SPAPIJsonParser
from spapi.models import AsinsInfo, SpapiPrices
from spapi.models import SpapiFees
from keepa.models import KeepaProducts
from keepa.models import convert_unix_time_to_keepa_time
from keepa.models import convert_recharts_data
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


class UpdateChartDataRequestTask(object):

    def __init__(self) ->  None:
        self.mq = MQ('chart')
        self.spapi_client = SPAPI()

    async def main(self):
        await self._get_chart_data_request()

    async def _get_chart_data_request(self, sleep_sec: int=60,
                                           limit_count: int=20,
                                           interval_sec: float=1.3):

        while True:
            asins = KeepaProducts.get_products_not_modified()
            if not asins:
                await asyncio.sleep(sleep_sec)
            
            for i in range(0, len(asins), limit_count):
                response = await self.spapi_client.get_competitive_pricing(asins[i:i+limit_count])
                self.mq.publish(json.dumps(response))
                await asyncio.sleep(interval_sec)


class UpdateChartData(object):

    def __init__(self) -> None:
        self.mq = MQ('chart')

    async def main(self):
        await self._update_chart_data_for_keepa_products()

    async def _update_chart_data_for_keepa_products(self, sleep_sec: int=60):

        for messages in self.mq.receive(100):
            if messages is None:
                await asyncio.sleep(sleep_sec)
                continue
            parsed_data = list(reduce(lambda data, func: map(func, data),
                [json.loads,
                 SPAPIJsonParser.parse_get_competitive_pricing], messages))
            parsed_data = ChainMap(*[{data['asin']: data for data in itertools.chain.from_iterable(parsed_data)}])
            keepa_products = await KeepaProducts.get_keepa_products_by_asins(parsed_data.keys())
            if keepa_products is None:
                continue
            with multiprocessing.Pool() as pool:
                keepa_products = pool.map(
                    partial(UpdateChartData._mapping_keepa_products_and_parsed_data,
                    parsed_data=parsed_data), keepa_products)
            await KeepaProducts.insert_all_on_conflict_do_update_chart_data(keepa_products)

    @staticmethod
    def _mapping_keepa_products_and_parsed_data(product: KeepaProducts, parsed_data: dict):
        now = convert_unix_time_to_keepa_time(time.time())    
        value = parsed_data.get(product.asin)
        if not value:
            return product

        price = value.get('price')
        rank = value.get('ranking')
        if not all((re.fullmatch('-?[0-9]+', str(price)), re.fullmatch('-?[0-9]+', str(rank)))):
            logger.error({'messagee': 'parameter is valid', "value": value})
            return product

        product.price_data[now] = price
        product.rank_data[now] = rank
        product.render_data = convert_recharts_data({
                                                'rank_data': product.rank_data,
                                                'price_data': product.price_data,})
        return product


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

        for messages in self.search_catalog_queue.receive():
            if messages is None:
                logger.info({'message': 'get_objects is None'})
                await asyncio.sleep(10)
                continue

            params = sorted([json.loads(resp) for resp in messages], key=lambda x: x['jan'])
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

            for i in range(0, len(asin_list), 20):
                response = await self.client.get_my_fees_estimates(asin_list[i:i+20])
                products = SPAPIJsonParser.parse_get_my_fees_estimates(response)
                for product in products:
                    result.append(SpapiFees(product['asin'], product['fee_rate'], product['ship_fee']))
                await asyncio.sleep(interval_sec)
            return result

        def _mws_mapping_spapi_fees(mws_list: List[MWS], spapi_fees: List[SpapiFees]) -> List[MWS]:
            mws_objects = deepcopy(mws_list)
            chain_map_fees = collections.ChainMap(
                                    *[{fee.asin: fee} for fee in spapi_fees])

            for mws in mws_objects:
                fee = chain_map_fees.get(mws.asin)
                if fee:
                    mws.fee_rate = fee.fee_rate
                    mws.shipping_fee = fee.ship_fee

            return mws_objects

        while True:
            mws_objects = await MWS.get_fee_is_None_asins(1000)
            if not mws_objects:
                await asyncio.sleep(30)
                continue

            asins = {mws.asin for mws in mws_objects}
            fees = await SpapiFees.get_asins_fee(list(asins))
            result = await _get_my_fees_estimate(list(asins - {fee.asin for fee in fees}))
            mws_objects = _mws_mapping_spapi_fees(mws_objects, fees + result)

            asyncio.ensure_future(SpapiFees.insert_all_on_conflict_do_update_fee(result))
            await MWS.insert_all_on_conflict_do_update_fee(mws_objects)
