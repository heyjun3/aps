import time
from multiprocessing import Process, Queue
import json
import asyncio

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

    def __init__(self, limit: int=20) -> None:
        self.queue = Queue()
        self.asins = KeepaProducts.get_products_not_modified()
        self.asins = [self.asins[i:i+limit] for i in range(0, len(self.asins), limit)]
        self.spapi_client = SPAPI()

    async def main(self) -> None:
        logger.info('action=main status=run')
        update_data_process = Process(target=self.update_data, args=(self.queue, ))
        update_data_process.start()
        get_competitive_pricing_task = asyncio.create_task(self.get_competitive_pricing())

        try:
            await asyncio.wait({get_competitive_pricing_task})
            update_data_process.join()
        except asyncio.TimeoutError as ex:
            logger.error(f'action=main error={ex}')
            self.queue.put(None)

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

        for asin_list in self.asins:
            response = await self.spapi_client.get_competitive_pricing(asin_list)
            products = SPAPIJsonParser.parse_get_competitive_pricing(response)
            [self.queue.put(product) for product in products]
            await asyncio.sleep(interval_sec)
        
        self.queue.put(None)
        logger.info('action=get_competitive_pricing status=done')


class RunAmzTask(object):

    def __init__(self, queue_name: str='mws', maxsize: int=1000) -> None:
        self.mq = MQ(queue_name)
        self.client = SPAPI()
        self.queue = asyncio.Queue(maxsize=maxsize)
        self.fees_queue = asyncio.Queue()

    async def main(self) -> None:
        logger.info('action=main status=run')

        get_mq_task = asyncio.create_task(self.get_mq())
        search_catalog_items_task = asyncio.create_task(self.search_catalog_items_v20220401())
        get_item_offers_task = asyncio.create_task(self.get_item_offers_batch())
        get_my_fee_estimate_task = asyncio.create_task(self.get_my_fees_estimate())
        await asyncio.wait({
            get_mq_task, 
            search_catalog_items_task,
            get_item_offers_task,
            get_my_fee_estimate_task,
        }, return_when='FIRST_COMPLETED')

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

            if not products:
                self.queue.put_nowait(params)
            else:
                for product in products:
                    MWS(asin=product['asin'], filename=params['filename'], title=product['title'],
                                jan=params['jan'], unit=product['quantity'], cost=params['cost']).save()

    async def search_catalog_items_v20220401(self, id_type: str='JAN', interval_sec: int=2) -> None:
        logger.info('action=search_catalog_items status=run')

        while True:
            params = [await self.queue.get() for _ in range(20) if not self.queue.empty()]
            logger.info({'queue_size': self.queue.qsize()})
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

    async def get_item_offers_batch(self, interval_sec: int=2):
        logger.info({'action': 'get_item_offers_batch', 'status': 'run'})

        async def _get_item_offers_batch(asins):
            logger.info({'action': '_get_item_offers_batch', 'status': 'run'})
            response = await self.client.get_item_offers_batch(asins)
            products = SPAPIJsonParser.parse_get_item_offers_batch(response)
            for product in products:
                MWS.update_price(asin=product['asin'], price=product['price'])
                SpapiPrices(asin=product['asin'], price=product['price'])
            
            logger.info({'action': '_get_item_offers_batch', 'status': 'done'})

        while True:
            asin_list = MWS.get_price_is_None_asins()
            if asin_list:
                asin_list = [asin_list[i:i+20] for i in range(0, len(asin_list), 20)]
                for asins in asin_list:
                    task = asyncio.create_task(_get_item_offers_batch(asins))
                    sleep = asyncio.create_task(asyncio.sleep(interval_sec))
                    await asyncio.gather(task, sleep)
            else:
                await asyncio.sleep(10)

    async def get_item_offers(self, interval_sec: int=0.2):
        logger.info('action=get_item_offers status=run')

        while True:
            asin_list = MWS.get_price_is_None_asins()
            if asin_list:
                for asin in asin_list:
                    response = await self.client.get_item_offers(asin)
                    product = SPAPIJsonParser.parse_get_item_offers(response)
                    MWS.update_price(asin=product['asin'], price=product['price'])
                    SpapiPrices(asin=product['asin'], price=product['price']).upsert()
                    await asyncio.sleep(interval_sec)
            else:
                await asyncio.sleep(10)

    async def get_my_fees_estimate(self, interval_sec: int=2, use_cache: bool=True) -> None:
        logger.info('action=get_my_fees_estimate_for_asin status=run')

        async def _get_fee_in_db(asin_list):
            if not use_cache:
                [self.fees_queue.put_nowait(asin) for asin in asin_list]
                return

            for asin in asin_list:
                fee = SpapiFees.get(asin)
                if fee is None:
                    self.fees_queue.put_nowait(asin)
                else:
                    MWS.update_fee(fee['asin'], fee['fee_rate'], fee['ship_fee'])

        async def _get_my_fees_estimate():
            async def _wrapper_get_my_fees_estimate():
                asins = [await self.fees_queue.get() for _ in range(20) if not self.fees_queue.empty()]
                response = await self.client.get_my_fees_estimates(asins)
                products = SPAPIJsonParser.parse_get_my_fees_estimates(response)
                for product in products:
                    SpapiFees(asin=product['asin'], fee_rate=product['fee_rate'], ship_fee=product['ship_fee']).upsert()
                    MWS.update_fee(asin=product['asin'], fee_rate=product['fee_rate'], shipping_fee=product['ship_fee'])

            while not self.fees_queue.empty():
                task = asyncio.create_task(_wrapper_get_my_fees_estimate())
                sleep = asyncio.create_task(asyncio.sleep(interval_sec))
                await asyncio.gather(task, sleep)

        while True:
            asin_list = MWS.get_fee_is_None_asins()
            if asin_list:
                get_fee_in_db_task = asyncio.create_task(_get_fee_in_db(asin_list))
                get_fees_estimate_task = asyncio.create_task(_get_my_fees_estimate())
                await asyncio.gather(get_fees_estimate_task, get_fee_in_db_task)
            else:
                await asyncio.sleep(30)

    async def get_my_fees_estimate_for_asin(self, interval_sec: int=0.2, use_cache: bool=True):
        logger.info('action=get_my_fees_estimate_for_asin status=run')

        async def _get_my_fees_estimate_for_asin(asin):
            response = await self.client.get_my_fees_estimate_for_asin(asin)
            product = SPAPIJsonParser.parse_get_my_fees_estimate_for_asin(response)
            try:
                SpapiFees(asin=product['asin'], fee_rate=product['fee_rate'], ship_fee=product['ship_fee']).upsert()
                MWS.update_fee(asin=product['asin'], fee_rate=product['fee_rate'], shipping_fee=product['ship_fee'])
            except KeyError as ex:
                logger.error(f'action=get_my_fees_estimate_for_asin error={ex}')

        while True:
            asin_list = MWS.get_fee_is_None_asins()
            if asin_list:
                for asin in asin_list:
                    if use_cache:
                        fee = SpapiFees.get(asin)
                    else:
                        fee = None
                    if fee is None:
                        await _get_my_fees_estimate_for_asin(asin)
                        await asyncio.sleep(interval_sec)
                    else:
                        try:
                            MWS.update_fee(asin=fee['asin'], fee_rate=fee['fee_rate'], shipping_fee=fee['ship_fee'])
                        except KeyError as ex:
                            logger.error(f'action=get_my_fees_estimate_for_asin error={ex}')
            else:
                await asyncio.sleep(10)

