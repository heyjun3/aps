from __future__ import annotations
import urllib.parse
import re
import datetime
import json
from pathlib import Path
from functools import partial
from typing import List
from typing import Callable
from dataclasses import dataclass
from functools import reduce

import gspread
import requests_html
from oauth2client.service_account import ServiceAccountCredentials

import log_settings
from mq import MQ
from crawler import utils
from crawler.buffalo.buffalo import BuffaloHTMLPage
from crawler.pc4u.pc4u import Pc4uHTMLPage
from crawler.pcones.pcones import PconesHTMLPage
from crawler.rakuten.rakuten import RakutenHTMLPage
from crawler.paypaymall.paypaymoll import PayPayMollHTMLParser


logger = log_settings.get_logger(__name__)

def log_decorator(func: Callable) -> Callable:
    def _inner(*args, **kwargs):
        logger.info({'action': func.__name__, 'status': 'run'})
        result = func(*args, **kwargs)
        logger.info({'action': func.__name__, 'status': 'done'})
        return result
    return _inner

@dataclass
class SpreadSheetValue:
    jan: str
    url: str


@dataclass
class RequestResult:
    jan: str
    url: str
    response: requests_html.HTMLResponse


@dataclass
class ParsedValue:
    jan: str
    price: int
    is_stocked: bool


@dataclass
class ParseResult:
    jan: str
    url: str
    parse_result: ParsedValue


class SpreadSheetCrawler(object):

    scope = ('https://spreadsheets.google.com/feeds', 'https://www.googleapis.com/auth/drive')
    requires = ('URL', 'JAN')

    def __init__(self, credential_file_name: str, sheet_title: str, sheet_name: str, queue_name: str='mws') -> None:
        self.credential_file_path = Path.cwd().joinpath(credential_file_name)
        self.credential = ServiceAccountCredentials.from_json_keyfile_name(self.credential_file_path, self.scope)
        self.client = gspread.authorize(self.credential)
        self.sheet_title = sheet_title
        self.sheet_name = sheet_name
        self.start_time = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')
        self.mq = MQ(queue_name)

    def start_crawler(self) -> None:
        sheet_values = self._get_crawl_urls_from_spread_sheet()
        # funcs = (
        #     self._validation_sheet_value,
        #     self._send_request,
        #     self._parse_response,
        #     self._generate_string_for_enqueue,
        #     self.mq.publish,
        # )
        # list(map(partial(self._request_sequence, funcs=funcs), sheet_values))
        publish_messages = filter(None, reduce(lambda d, f: map(f, d), [
            self._validation_sheet_value,
            self._send_request,
            self._parse_response,
            self._generate_string_for_enqueue,
        ], sheet_values))
        [self.mq.publish(message) for message in publish_messages]

    # def _request_sequence(self, value: dict, funcs: tuple[Callable]) -> dict|None:
    #     if value is None:
    #         return
    #     if not funcs:
    #         return value

    #     return self._request_sequence(funcs[0](value), funcs[1:])

    def _get_crawl_urls_from_spread_sheet(self) -> List[SpreadSheetValue]:
        sheet = self.client.open(self.sheet_title).worksheet(self.sheet_name)
        records = list(map(lambda x: SpreadSheetValue(x.get('JAN'), x.get('URL')), sheet.get_all_records()))
        return records

    def _validation_sheet_value(self, value: SpreadSheetValue) -> SpreadSheetValue|None:
        if value is None: return
        if value.url is None:
            logger.error({'message': 'sheet value is URL None'})
            return

        return value

    def _send_request(self, sheet_value: SpreadSheetValue, interval_sec: int=4) -> RequestResult|None:
        if sheet_value is None: return
        logger.info(sheet_value.url)
        response = utils.request(sheet_value.url, time_sleep=interval_sec)

        if response is None:
            return

        return RequestResult(sheet_value.jan, sheet_value.url, response)
        # if response.status_code == 200:
            # return RequestResult(sheet_value.jan, sheet_value.url, response)
        # if response.status_code == 404:
        #     logger.error({'status_code': response.status_code, 'message': 'page not Found'})
        #     return
        # logger.error(response.status_code, response.url)

    def _parse_response(self, res: RequestResult) -> ParseResult|None:
        if res is None: return
        parser = self._get_html_parser(res.response.url)

        if parser is None:
            return
        value = parser(res.response.text)
        parsed_value = ParsedValue(value.get('jan'), value.get('price'), value.get('is_stocked'))
        return ParseResult(res.jan, res.url, parsed_value)

    def _generate_string_for_enqueue(self, result: ParseResult) -> str|None:
        if result is None: return

        jan = jan if (jan := result.parse_result.jan) else result.jan
        price = result.parse_result.price
        is_stocked = result.parse_result.is_stocked
        if not all((jan, price, is_stocked)):
            logger.error({'message': 'publish queue bad parameter', 'values': (jan, price, is_stocked)})
            return

        return json.dumps({
            'filename': f'repeat_{self.start_time}',
            'jan': jan,
            'cost': price,
            'url': result.url,
        })

    def _get_html_parser(self, url: str) -> Callable:
        netloc = urllib.parse.urlparse(url).netloc
        # todo add parser
        if re.search('(geno-web.jp)$', netloc):
            return
        if re.search('(janpara.co.jp)$', netloc):
            return
        if re.search('(system5.jp)$', netloc):
            return
        if re.search('(pc-koubou.jp)$', netloc):
            return
        if re.search('(netmall.hardoff.co.jp)$', netloc):
            return
        if re.search('(item.rakuten.co.jp)$', netloc):
            return RakutenHTMLPage.scrape_product_detail_page
        if re.search('(pc4u.co.jp)$', netloc):
            return Pc4uHTMLPage.scrape_product_detail_page
        if re.search('(1-s.jp)$', netloc):
            return PconesHTMLPage.scrape_product_detail_page
        if re.search('(buffalo-direct.com)$', netloc):
            return BuffaloHTMLPage.scrape_product_detail_page
        if re.search('(paypaymall.yahoo.co.jp)$', netloc):
            return PayPayMollHTMLParser.product_detail_page_parser
        if re.search('(sofmap.com)$', netloc):
            return
        if re.search('(soundhouse.co.jp)$', netloc):
            return
        if re.search('(ec.treasure-f.com)$', netloc):
            return
        if re.search('(e-trend.co.jp)$', netloc):
            return
        if re.search('(ikebe-gakki.com)$', netloc):
            return
        if re.search('(pioneer-itstore.jp)$', netloc):
            return

        logger.error({'message': 'netloc is not match', 'value': netloc})
