from __future__ import annotations
import re
import json
import functools
from functools import partial
from dataclasses import dataclass
from dataclasses import asdict
from typing import List
from copy import deepcopy
from datetime import datetime

from bs4 import BeautifulSoup
from requests_html import HTMLResponse

import log_settings
import mq
from crawler import utils


logger = log_settings.get_logger(__name__)
logging = log_settings.decorator_logging(logger)


class YahooShopApi(object):

    @staticmethod
    def item_search_v3(request: ItemSearchRequest, interval_sec=2) -> HTMLResponse:
        endpoint = 'https://shopping.yahooapis.jp/ShoppingWebService/V3/itemSearch'
        res = utils.request(endpoint, params=asdict(request), time_sleep=interval_sec)
        return res

@dataclass
class ItemSearchRequest:
    appid: str
    seller_id: str
    condition: str = 'new'
    in_stock: str = 'true'
    price_to: int = 100000
    results: int = 100
    sort: str = '-price'
    start: int = 1

class YahooShopApiParser(object):

    @staticmethod
    def parse_item_search_v3(response: dict) -> List[ItemSearchResult]:
        result = []
        for item in response.get('hits'):
            match item:
                case {
                    'code': code, 'price': price, 'janCode': jan, 'name': name,
                    'point': {'premiumBonusAmount': point},
                    'seller': {'sellerId': sellerId},
                    'url': url, }:
                    result.append(ItemSearchResult(code, price, jan, name, point, sellerId, url))

        return result

@dataclass
class ItemSearchResult:
    product_id: str
    price: int
    jan: str = None
    name: str = None
    point: int = None
    shop_id: str = None
    url: str = None


class YahooShopCrawler(object):
    def __init__(self):
        pass
    
    @logging
    def search_by_shop_id(self, app_id: str, seller_id: str, mq: mq.MQ=mq.MQ('mws')) -> None:
        timestamp = datetime.now()
        query = ItemSearchRequest(app_id, seller_id)
        while query:
            logger.info(query)
            res = YahooShopApi.item_search_v3(query)
            messages = self._search_sequence(res.json(), timestamp)
            [mq.publish(message) for message in messages if message]
            query = self._generate_next_query(res.json(), query)

    @logging
    def _search_sequence(self, res: dict, timestamp: datetime) -> map:
        results = functools.reduce(lambda d, f: f(d), [
            YahooShopApiParser.parse_item_search_v3,
            partial(map, self._calc_real_price),
            partial(map, partial(self._generate_publish_message, timestamp=timestamp)),
        ], res)
        return results

    def _calc_real_price(self, item: ItemSearchResult) -> ItemSearchResult|None:
        result = deepcopy(item)
        match result:
            case ItemSearchResult(price=price, point=point) if price and point is not None:
                result.price = price - point
                return result
            case _ :
                logger.error({
                    "message": "invalid value",
                    "action": "_calc_real_price",
                    "value": result})
                return

    def _generate_publish_message(self, item: ItemSearchResult,
                            timestamp: datetime, prefix: str='paypay') -> str|None:
        match item, timestamp:
            case ItemSearchResult(jan=jan, price=price, url=url), datetime() if all((jan, price, url)):
                return json.dumps({
                    'jan': jan, 'cost': price, 'url': url,
                    "filename": f'{prefix}_{timestamp.strftime("%Y%m%d_%H%M%S")}'})
            case _ :
                logger.error({
                    "message": "invalid value",
                    "action": "_generate_publish_message",
                    "value": item})
                return

    def _generate_next_query(self, res: dict, query: ItemSearchRequest) -> ItemSearchRequest:
        request = deepcopy(query)
        match res:
            case {"firstResultsPosition": 900, "totalResultsReturned": 100,
                  "hits": [*_, {"price": last_item_price}]}:
                request.price_to = last_item_price - 1
                request.start = 1
                return request
            case {"firstResultsPosition": 1, "totalResultsReturned": 100}:
                request.start = 100
                return request
            case {"firstResultsPosition": position, "totalResultsReturned": 100} if position <= 900:
                request.start = position + 100
                return request
            case _:
                return


@dataclass
class ParsedPayPayMollDetailPage:
    jan: str
    price: int
    is_stocked: bool

class PayPayMollHTMLParser(object):

    @staticmethod
    def product_detail_page_parser(response: str) -> dict:

        soup = BeautifulSoup(response, 'lxml')
        try:
            item_details = soup.select('.ItemDetails_list')
            item_detail = list(filter(None, map(lambda x: ''.join(re.findall('[0-9]', x.text)), item_details)))
            jan = item_detail.pop() if item_detail else None
        except (AttributeError) as ex:
            logger.info({'message': 'jan is None', 'error': ex})
            jan = None

        price = soup.select_one('.ItemPrice_price')
        if price:
            price = int(''.join(re.findall('[0-9]', price.text)))

        is_stocked = soup.select_one('#CartButtonUltLog').attrs
        if 'disabled' in is_stocked:
            is_stocked = False

        return {'jan': jan, 'price': price, 'is_stocked': bool(is_stocked)}
