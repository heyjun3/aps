import re
from dataclasses import dataclass

from bs4 import BeautifulSoup
from requests_html import HTMLResponse

import log_settings
from crawler import utils


logger = log_settings.get_logger(__name__)


class YahooShoppingApiClient(object):

    def __init__(self):
        pass

    def item_search_v3(self, query: dict, interval_sec=1) -> HTMLResponse:
        endpoint = 'https://shopping.yahooapis.jp/ShoppingWebService/V3/itemSearch'
        res = utils.request(endpoint, params=query, time_sleep=interval_sec)
        return res


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
