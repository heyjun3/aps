import os
import unittest
from unittest.mock import MagicMock

import pytest
from crawler.rakuten.models import RakutenProduct

from crawler.rakuten.rakuten import RakutenHTMLPage
from crawler.rakuten.rakuten import RakutenCrawler
from crawler.rakuten.rakuten import MaxProductsCountNotFoundException


dirname = os.path.join(os.path.dirname(__file__), 'test_html')


class TestRakutenCrawler(object):

    def test_generate_querys(self):
        client = RakutenCrawler('test', 'test')
        path = os.path.join(dirname, 'product_list_page.html')
        response = MagicMock()
        with open(path, 'r') as f:
            response.text = f.read()

        querys = client._generate_querys(response)
        assert querys[0]['p'] == 1
        assert querys[0]["s"] == 3
        assert querys[0]['used'] == 0
        assert querys[0]['sid'] == 'test'

        assert querys[-1]['p'] == 16
        assert querys[-1]["s"] == 3
        assert querys[-1]['used'] == 0
        assert querys[-1]['sid'] == 'test'

        assert len(querys) == 16

    def test_mapping_rakuten_products(self):
        values = [{'product_code': 'aaa', 'price': 111},
                  {'product_code': 'bbb', 'price': 222},]
        rakuten_products = [RakutenProduct(product_code='aaa', jan='9999'),
                            RakutenProduct(product_code='ccc', jan='0000')]
        client = RakutenCrawler('test', 'test')
        products = client._mapping_rakuten_products(values, rakuten_products)
        assert len(products) == len(values)
        assert products[0]['product_code'] == 'aaa'
        assert products[0]['price'] == 111
        assert products[0]['jan'] == '9999'
        assert products[-1]['product_code'] == 'bbb'
        assert products[-1]['price'] == 222
        assert 'jan' not in products[-1]


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
        
    def test_parse_max_products_count(self):
        path = os.path.join(dirname, 'product_list_page.html')
        with open(path, 'r') as f:
            response = f.read()
        count = RakutenHTMLPage.parse_max_products_count(response)
        assert count == 691

    def test_parse_max_products_count_fail(self):
        with pytest.raises(MaxProductsCountNotFoundException) as e:
            count = RakutenHTMLPage.parse_max_products_count('')

        assert str(e.value) == ''
