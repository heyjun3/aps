from __future__ import annotations
import re
import json
from dataclasses import dataclass
from dataclasses import asdict
from typing import List
from copy import deepcopy
from datetime import datetime

from bs4 import BeautifulSoup
from requests_html import HTMLResponse

import log_settings
from crawler import utils


logger = log_settings.get_logger(__name__)


class YahooShopApi(object):

    @staticmethod
    def item_search_v3(request: ItemSearchRequest, interval_sec=1) -> HTMLResponse:
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

    def search_by_shop_id(self, app_id: str, seller_id: str) -> None:
        query = ItemSearchRequest(app_id, seller_id)
        res = YahooShopApi.item_search_v3(query)
        results = YahooShopApiParser.parse_item_search_v3(res.json())
        values = [self._calc_real_price(result) for result in results]


    def _calc_real_price(self, item: ItemSearchResult) -> ItemSearchResult|None:
        result = deepcopy(item)
        match result:
            case ItemSearchResult(price=price, point=point) if all((price, point)):
                result.price = price - point
                return result
            case _ :
                return

    def _generate_publish_message(self, item: ItemSearchResult,
                            timestamp: datetime, prefix: str='paypay') -> str|None:
        match item, timestamp:
            case ItemSearchResult(jan=jan, price=price, url=url), datetime() if all((jan, price, url)):
                return json.dumps({
                    'jan': jan, 'cost': price, 'url': url,
                    "filename": f'{prefix}_{timestamp.strftime("%Y%m%d_%H%M%S")}'})
            case _ :
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
