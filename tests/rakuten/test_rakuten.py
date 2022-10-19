import unittest
import os

from crawler.rakuten.rakuten import RakutenHTMLPage
from crawler.rakuten.rakuten import RakutenCrawler


dirname = os.path.join(os.path.dirname(__file__), 'test_html')


class TestRakutenCrawler(object):

    def test_create_querys_success(self):
        client = RakutenCrawler('test', 'test')
        count = 9
        querys = client._create_querys(9)
        assert querys[0]['p'] == 1
        assert querys[0]["s"] == 3
        assert querys[0]['used'] == 0
        assert querys[0]['sid'] == 'test'

        assert querys[-1]['p'] == 9
        assert querys[-1]["s"] == 3
        assert querys[-1]['used'] == 0
        assert querys[-1]['sid'] == 'test'

        assert len(querys) == count


class ScrapeDetailProductPage(unittest.TestCase):

    def test_scrape_jan_code(self):
        
        html_path = os.path.join(dirname, 'scrape_jan_code.html')
        with open(html_path, 'r') as f:
            response = f.read()

        parsed_value = RakutenHTMLPage.scrape_product_detail_page(response)
        self.assertEqual(parsed_value.get('jan'), '4589919807796')
        self.assertEqual(parsed_value.get('price'),16000) 
        self.assertEqual(parsed_value.get('is_stocked'), True)
        # self.assertEqual(parsed_value.get('point'), 2555)
        
    def test_parse_product_list_page(self):

        path = os.path.join(dirname, 'product_list_page.html')
        with open(path, 'r') as f:
            response = f.read()

        parsed_value = RakutenHTMLPage.parse_product_list_page(response)
        assert parsed_value[0]['name'] == 'シャープ 加湿空気清浄機 KI-NS40W'
        assert parsed_value[0]['url'] == 'https://item.rakuten.co.jp/superdeal/11118kins40w20211101/'
        assert parsed_value[0]['price'] == 29800
        assert parsed_value[0]['product_code'] == '11118kins40w20211101'

        assert parsed_value[-1]['name'] == 'シロカ siroca 4L 電気圧力鍋 SP-4D151'
        assert parsed_value[-1]['url'] == 'https://item.rakuten.co.jp/superdeal/10052sp-4d151191004/'
        assert parsed_value[-1]['price'] == 16000
        assert parsed_value[-1]['product_code'] == '10052sp-4d151191004'
        