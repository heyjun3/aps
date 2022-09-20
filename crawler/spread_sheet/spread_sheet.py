import time
import urllib.parse
import re
import datetime
import json
from pathlib import Path
from functools import partial
from copy import deepcopy
from typing import List
from typing import Callable

import gspread
import requests
from oauth2client.service_account import ServiceAccountCredentials

import log_settings
from mq import MQ
from crawler.utils import HEADERS
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

    @log_decorator
    def start_crawler(self) -> None:
        sheet_values = self._get_crawl_urls_from_spread_sheet()
        funcs = (
            self._validation_sheet_value,
            self._send_request,
            self._parse_response,
            self._publish_queue,
        )
        list(map(partial(self._request_sequence, funcs=funcs), sheet_values))

    @log_decorator
    def _request_sequence(self, value: dict, funcs: tuple[Callable]) -> dict|None:
        if value is None:
            return
        if not funcs:
            return value

        return self._request_sequence(funcs[0](value), funcs[1:])

    @log_decorator
    def _get_crawl_urls_from_spread_sheet(self) -> List[dict]:
        sheet = self.client.open(self.sheet_title).worksheet(self.sheet_name)
        return sheet.get_all_records()

    @log_decorator
    def _validation_sheet_value(self, sheet_value: dict) -> dict|None:
        value = deepcopy(sheet_value)
        if not all(require in value for require in self.requires):
            logger.error({'message': 'sheet_value is not valid headers', 'value': sheet_value})
            return

        if value.get('URL') is None:
            logger.error({'message': 'sheet value is URL None'})
            return

        return value

    @log_decorator
    def _send_request(self, sheet_value: dict, interval_sec: int=4) -> dict|None:
        value = deepcopy(sheet_value)
        url = value.get('URL')
        logger.info(url)

        response = requests.get(url, headers=HEADERS)
        time.sleep(interval_sec)

        if response.status_code == 200:
            value['response'] = response
            return value
        if response.status_code == 404:
            logger.error({'status_code': response.status_code, 'message': 'page not Found'})
            return
        logger.error(response.status_code, response.url)

    @log_decorator
    def _parse_response(self, response: dict) -> dict|None:
        value = deepcopy(response)
        response = value.get('response')
        parser = self._get_html_parser(response.url)

        if parser is None:
            return
        value['parse_value'] = parser(response.text)
        return value

    @log_decorator
    def _publish_queue(self, response: dict) -> None:

        jan = response['parse_value'].get('jan') or response.get('JAN')
        price = response['parse_value'].get('price')
        is_stocked = response['parse_value'].get('is_stocked')
        if not all((jan, price, is_stocked)):
            logger.error({'message': 'publish queue bad parameter', 'values': (jan, price, is_stocked)})
            return

        self.mq.publish(json.dumps({
            'filename': f'repeat_{self.start_time}',
            'jan': jan,
            'cost': price,
            'url': response['response'].url,
        }))

    @log_decorator
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
