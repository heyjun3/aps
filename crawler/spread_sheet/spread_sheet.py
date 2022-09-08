from pathlib import Path

import gspread
from oauth2client.service_account import ServiceAccountCredentials


class SpreadSheetCrawler(object):

    scope =['https://spreadsheets.google.com/feeds', 'https://www.googleapis.com/auth/drive']

    def __init__(self, credential_file_name: str) -> None:
        self.credential_file_path = Path.cwd().joinpath(credential_file_name)
        self.credential = ServiceAccountCredentials.from_json_keyfile_name(self.credential_file_path, self.scope)
        self.client = gspread.authorize(self.credential)

    def get_crawl_urls_from_spread_sheet(self):
        sheet = self.client.open('bussiness').worksheet('repeat_list')
        print(sheet.get_all_records())


if __name__ == '__main__':
    client = SpreadSheetCrawler('gsheet-355401-5fbc168f98c2.json')
    client.get_crawl_urls_from_spread_sheet()
