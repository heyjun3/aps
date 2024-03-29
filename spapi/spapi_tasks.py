from copy import deepcopy
import multiprocessing
import time
import json
import itertools
import asyncio
import datetime
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
from keepa import convert
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
        self.price_queue = MQ("price")
        self.spapi_client = SPAPI()

    async def main(self):
        await self._get_chart_data_request()

    async def _get_chart_data_request(self, sleep_sec: int=60,
                                           limit_count: int=20,
                                           interval_sec: float=1.3):

        while True:
            asins = KeepaProducts.get_products_not_modified()
            if not asins:
                asins_price = list(filter(None, [self.price_queue.basic_get() for _ in range(20)]))
                if not asins_price:
                    await asyncio.sleep(sleep_sec)
                    continue
                res = await self.spapi_client.get_competitive_pricing(asins_price)
                products = SPAPIJsonParser.parse_get_competitive_pricing(res)
                asyncio.ensure_future(MWS.bulk_update_prices(products))
                asyncio.ensure_future(SpapiPrices.insert_all_on_conflict_do_update_price(products))
                await asyncio.sleep(interval_sec)
                continue
            
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
        logger.info({"action": "mapping_keepa_products_and_parsed_data", "status": "run"})
        now = convert.unix_time_to_keepa_time(time.time())    
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
        if product.render_data is None:
            logger.error({
                "action": "mapping_keepa_products_and_parsed_data",
                "message": f'render_data is null. asin:{product.asin}'},)
            return product
        
        chart_data = product.render_data.get("data")
        if chart_data is None:
            return product
        
        today = datetime.datetime.now().date().strftime("%Y-%m-%d")
        chart_data = list(filter(UpdateChartData._date_filter, chart_data))
        chart_data.append({"date": today, "rank": rank, "price": price})
        product.render_data["data"] = chart_data

        logger.info({"action": "mapping_keepa_products_and_parsed_data", "status": "done"})
        return product
    
    @staticmethod
    def _date_filter(value: dict, past_days: int=90) -> bool:
        date = value.get("date")
        if date is None:
            return False
        try:
            date = datetime.datetime.strptime(date, "%Y-%m-%d")
        except ValueError as e:
            logger.error({"message": "failed parse date", "date": date, "error": e})
            return False
        past_date = datetime.datetime.now() - datetime.timedelta(days=past_days)
        return past_date < date


class RunAmzTask(object):

    def __init__(self, queue_name: str='mws', search_queue: str='search_catalog', price_queue: str='price') -> None:
        self.mq = MQ(queue_name)
        self.search_catalog_queue = MQ(search_queue)
        self.price_queue = MQ(price_queue)
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

    @log_decorator
    def _validation_parameter(self, param: dict) -> dict|None:
        require = ('cost', 'jan', 'filename', 'url')
        if not all(r in param for r in require):
            logger.error({'bad parameter': param})
            return
        return param

    @log_decorator
    def _map_asin_info_and_message(self, messages: List[dict], asins: List[AsinsInfo]):
        send_messages, mws_objects = [], []
        asin_infos = {k: list(g) for k, g in itertools.groupby(
                            sorted(asins, key=lambda x: x.jan), lambda x: x.jan)}

        for message in messages:
            asin_info = asin_infos.get(message['jan'])
            if not asin_info:
                send_messages.append(message)
                continue
            mws_object = [MWS(
                            asin=info.asin,
                            filename=message['filename'],
                            title=info.title,
                            jan=info.jan,
                            unit=info.quantity,
                            cost=message['cost'],
                            url=message['url'],
                            ) for info in asin_info]
            mws_objects.extend(mws_object)

        return send_messages, mws_objects

    async def get_queue(self, interval_sec: int=10, task_count=100) -> None:
        logger.info('action=get_queue status=run')

        for messages in self.mq.receive(task_count):
            if messages is None:
                await asyncio.sleep(interval_sec)
                continue

            messages = list(filter(None, reduce(lambda d, f: map(f, d), [
                json.loads, self._validation_parameter], messages)))
            asins = await AsinsInfo.get_asin_object_by_jan_list(
                [message.get('jan') for message in messages])

            if not asins:
                [self.search_catalog_queue.publish(json.dumps(message))
                                                         for message in messages]
                continue

            send_messages, mws_objects = self._map_asin_info_and_message(messages, asins)
            if mws_objects:
                await MWS.insert_all_on_conflict_do_nothing(mws_objects)
            if send_messages:
                [self.search_catalog_queue.publish(json.dumps(message)) for message in send_messages]

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

    async def get_item_offers_batch(self, interval_sec: int=2, count: int=20):
        logger.info({'action': 'get_item_offers_batch', 'status': 'run'})
                       
        while True:
            asins = await MWS.get_asins_by_price_is_None()
            if not asins:
                await asyncio.sleep(10)
                continue
            if not self.price_queue.get_message_count():
                [self.price_queue.publish(asin) for asin in asins]
            for messages in self.price_queue.receive(count):
                if not messages:
                    break
                resp = await self.client.get_item_offers_batch(messages)
                products = SPAPIJsonParser.parse_get_item_offers_batch(resp)
                asyncio.ensure_future(MWS.bulk_update_prices(products))
                asyncio.ensure_future(SpapiPrices.insert_all_on_conflict_do_update_price(products))
                await asyncio.sleep(interval_sec)

    async def get_item_offer(self, interval_sec: float=0.5):
        logger.info({"action": "get_item_offer", "status": "run"})

        for message in self.price_queue.get():
            if message is None:
                logger.info({"message": "price queue is empty"})
                time.sleep(10)
                continue
            
            res = await self.client.get_item_offers(message)
            product = SPAPIJsonParser.parse_get_item_offers(res)
            asyncio.ensure_future(MWS.bulk_update_prices([product]))
            asyncio.ensure_future(SpapiPrices.insert_all_on_conflict_do_update_price([product]))
            await asyncio.sleep(interval_sec)

    async def get_my_fees_estimate(self, interval_sec: int=2) -> None:
        logger.info('action=get_my_fees_estimate_for_asin status=run')

        while True:
            asins = await MWS.get_asins_by_fee_is_None(1000)
            if not asins:
                await asyncio.sleep(30)
                continue

            fees = await SpapiFees.get_asins_fee(asins)
            if fees:
                asyncio.ensure_future(MWS.bulk_update_fees([fee.values for fee in fees]))

            asins = list(set(asins) - {fee.asin for fee in fees})
            for i in range(0, len(asins), 20):
                resp = await self.client.get_my_fees_estimates(asins[i:i+20])
                products = SPAPIJsonParser.parse_get_my_fees_estimates(resp)
                await asyncio.gather(
                    MWS.bulk_update_fees(deepcopy(products)),
                    SpapiFees.insert_all_on_conflict_do_update_fee(products),
                    asyncio.sleep(interval_sec),
                )
