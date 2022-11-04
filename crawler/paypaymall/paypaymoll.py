from __future__ import annotations
import re
from dataclasses import dataclass
from dataclasses import asdict
from typing import List

from bs4 import BeautifulSoup
from requests_html import HTMLResponse

import log_settings
from crawler import utils


logger = log_settings.get_logger(__name__)


class YahooShoppingApiClient(object):

    @staticmethod
    def item_search_v3(request: YahooShoppingApiItemSearchRequest, interval_sec=1) -> HTMLResponse:
        endpoint = 'https://shopping.yahooapis.jp/ShoppingWebService/V3/itemSearch'
        res = utils.request(endpoint, params=asdict(request), time_sleep=interval_sec)
        return res

@dataclass
class YahooShoppingApiItemSearchRequest:
    appid: str
    seller_id: str
    condition: str = 'new'
    in_stock: str = 'true'
    results: int = 100
    sort: str = '-price'
    start: int = 1

class YahooShoppingApiParser(object):

    @staticmethod
    def parse_item_search_v3(response: dict) -> List[YahooShoppingSearchItem]:
        result = []
        for item in response.get('hits'):
            result.append(YahooShoppingSearchItem(
                product_id=item.get('code'),
                price=item.get('price'),
                jan=item.get('janCode'),
                name=item.get('name'),
                point=point.get('premiumAmount') if (point := item.get('point')) else None,
                shop_id=sid.get('sellerId') if (sid := item.get('serller')) else None,
                url=item.get('url'),
            ))
        return result

@dataclass
class YahooShoppingSearchItem:
    product_id: str
    price: int
    jan: str = None
    name: str = None
    point: int = None
    shop_id: str = None
    url: str = None


class YahooShoppingCrawler(object):
    def __init__(self):
        pass


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
