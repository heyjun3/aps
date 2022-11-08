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
            result.append(ItemSearchResult(
                product_id=item.get('code'),
                price=int(price) if (price := item.get('price')) else None,
                jan=item.get('janCode'),
                name=item.get('name'),
                point=int(point) if (point := item['point']['premiumBonusAmount']) else None,
                shop_id=sid.get('sellerId') if (sid := item.get('seller')) else None,
                url=item.get('url'),
            ))
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
