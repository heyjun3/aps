import time
import urllib.parse
from pathlib import Path
from copy import deepcopy
from typing import List
from typing import Callable

import gspread
import requests
from oauth2client.service_account import ServiceAccountCredentials

import log_settings


logger = log_settings.get_logger(__name__)


class SpreadSheetCrawler(object):

    scope = ('https://spreadsheets.google.com/feeds', 'https://www.googleapis.com/auth/drive')
    requires = ('URL', 'JAN')

    def __init__(self, credential_file_name: str, sheet_title: str, sheet_name: str) -> None:
        self.credential_file_path = Path.cwd().joinpath(credential_file_name)
        self.credential = ServiceAccountCredentials.from_json_keyfile_name(self.credential_file_path, self.scope)
        self.client = gspread.authorize(self.credential)
        self.sheet_title = sheet_title
        self.sheet_name = sheet_name

    def _get_crawl_urls_from_spread_sheet(self) -> List[dict]:
        sheet = self.client.open(self.sheet_title).worksheet(self.sheet_name)
        return sheet.get_all_records()

    def _validation_sheet_value(self, sheet_value: dict) -> dict|None:
        value = deepcopy(sheet_value)
        if value is None:
            return

        if not all(require in value for require in self.requires):
            return

        if not value.get('URL'):
            return

        return value

    def _send_request(self, sheet_value: dict, interval_sec: int=2) -> requests.Response|None:
        if sheet_value is None:
            return
        response = requests.get(sheet_value.get('URL'))
        time.sleep(interval_sec)

        if response.status_code == 200:
            return response
        if response.status_code == 404:
            return
        logger.errror(response.status_code, response.json())
        return

    def _get_html_parser(self, url: str) -> Callable:
        netloc = urllib.parse.urlparse(url).netloc
        if netloc == 'rakuten':
            return

        return 
        

    def _parse_response(self, response: requests.Response) -> dict|None:
        if response is None:
            return
        
        parser = self._get_html_parser(response.url)
        if parser is None:
            return


        


        


if __name__ == '__main__':
    client = SpreadSheetCrawler('gsheet-355401-5fbc168f98c2.json')
    client.get_crawl_urls_from_spread_sheet()
