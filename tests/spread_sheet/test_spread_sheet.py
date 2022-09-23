from crawler.spread_sheet.spread_sheet import SpreadSheetCrawler
from crawler.spread_sheet.spread_sheet import SpreadSheetValue
import settings

class TestSpreadSheet(object):

    def test_validation_sheet_value_success(self):
        crawler = SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, 'test', 'test')
        value = SpreadSheetValue('URL', 'JAN')
        assert crawler._validation_sheet_value(value) == value
