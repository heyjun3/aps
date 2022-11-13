from crawler.buffalo.buffalo import BuffaloHTMLPage
from crawler.paypaymall.paypaymoll import PayPayMollHTMLParser
from crawler.pc4u.pc4u import Pc4uHTMLPage
from crawler.pcones.pcones import PconesHTMLPage
from crawler.spread_sheet.spread_sheet import SpreadSheetCrawler
from crawler.spread_sheet.spread_sheet import SpreadSheetValue
from crawler.spread_sheet.spread_sheet import ParseResult
from crawler.spread_sheet.spread_sheet import ParsedValue
from crawler.rakuten.rakuten import RakutenHTMLPage
import settings

class TestSpreadSheet(object):

    def test_validation_sheet_value_success(self):
        crawler = SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, 'test', 'test')
        value = SpreadSheetValue('URL', 'JAN')
        assert crawler._validation_sheet_value(value) == value

    def test_validation_sheet_value_faild(self):
        crawler = SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, 'test', 'test')
        value = SpreadSheetValue('jan', None)
        assert crawler._validation_sheet_value(value) == None

    def test_generate_string_for_enqueue_success(self):
        crawler = SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, 'test', 'test')
        crawler.start_time ="20220923_155822"
        value = ParsedValue('jan', 1111, True)
        result = ParseResult('jan', 'URL', value)
        assert crawler._generate_string_for_enqueue(result) == '{"filename": "repeat_20220923_155822", "jan": "jan", "cost": 1111, "url": "URL"}'

    def test_generate_string_for_enqueue_fail_jan_is_None(self):
        crawler = SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, 'test', 'test')
        crawler.start_time ="20220923_155822"
        value = ParsedValue(None, 1111, True)
        result = ParseResult(None, 'URL', value)
        assert crawler._generate_string_for_enqueue(result) == None

    def test_generate_string_for_enqueue_fail_price_is_None(self):
        crawler = SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, 'test', 'test')
        crawler.start_time ="20220923_155822"
        value = ParsedValue('jan', None, True)
        result = ParseResult('jan', 'URL', value)
        assert crawler._generate_string_for_enqueue(result) == None

    def test_generate_string_for_enqueue_fail_is_stocked_false(self):
        crawler = SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, 'test', 'test')
        crawler.start_time ="20220923_155822"
        value = ParsedValue('jan', 1111, False)
        result = ParseResult('jan', 'URL', value)
        assert crawler._generate_string_for_enqueue(result) == None

    def test_get_html_parser_success(self):
        crawler = SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, 'test', 'test')
        rakuten_url = 'https://item.rakuten.co.jp'
        pc4u_url = 'https://pc4u.co.jp'
        ones_url = 'https://1-s.jp'
        buffalo_url = 'https://buffalo-direct.com'
        paypaymall_url = 'https://paypaymall.yahoo.co.jp'

        assert crawler._get_html_parser(rakuten_url) == RakutenHTMLPage.scrape_product_detail_page
        assert crawler._get_html_parser(pc4u_url) == Pc4uHTMLPage.scrape_product_detail_page
        assert crawler._get_html_parser(ones_url) == PconesHTMLPage.scrape_product_detail_page
        assert crawler._get_html_parser(buffalo_url) == BuffaloHTMLPage.scrape_product_detail_page
        assert crawler._get_html_parser(paypaymall_url) == PayPayMollHTMLParser.product_detail_page_parser

    def test_get_html_parser_faild(self):
        crawler = SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, 'test', 'test')
        url = 'https://google.com'

        assert crawler._get_html_parser(url) == None
