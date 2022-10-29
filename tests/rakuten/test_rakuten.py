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

    def test_mapping_rakuten_products(self):
        values = [{'product_code': 'aaa', 'price': 111},
                  {'product_code': 'bbb', 'price': 222},]
        rakuten_products = [RakutenProduct(product_code='aaa', jan='9999'),
                            RakutenProduct(product_code='ccc', jan='0000')]
        client = RakutenCrawler()
        products = client._mapping_rakuten_products(values, rakuten_products)
        assert len(products) == len(values)
        assert products[0]['product_code'] == 'aaa'
        assert products[0]['price'] == 111
        assert products[0]['jan'] == '9999'
        assert products[-1]['product_code'] == 'bbb'
        assert products[-1]['price'] == 222
        assert 'jan' not in products[-1]

    def test_mapping_rakuten_products_2(self):
        values = [{'product_code': '4000000000000', 'price': 1000}]
        rakuten_products = [RakutenProduct(product_code='aaa', jan='1111')]
        client = RakutenCrawler()
        products = client._mapping_rakuten_products(values, rakuten_products)
        assert products[0]['product_code'] == '4000000000000'
        assert products[0]['price'] == 1000

    def test_mapping_search_value(self):
        values = [{'product_code': 'aaa', 'price': 111},
                  {'product_code': 'bbb', 'price': 222},]
        value = {'product_code': 'aaa', 'jan': '0000'}
        client = RakutenCrawler()
        product = client._mapping_search_value(value, values)
        assert product.get('product_code') == 'aaa'
        assert product.get('jan') == '0000'
        assert product.get('price') == 111

    def test_mapping_search_value_None(self):
        values = [{'product_code': 'aaa', 'price': 111},
                  {'product_code': 'bbb', 'price': 222},]
        value = {'product_code': 'ccc', 'jan': '0000'}
        client = RakutenCrawler()
        product = client._mapping_search_value(value, values)
        assert product == None

    def test_calc_real_price(self):
        value = {'product_code': 'aaa', 'price': 10000, 'point': 1000}
        client = RakutenCrawler()
        result = client._calc_real_price(value)
        assert result.get('price') == 8000
        assert result.get('product_code') == 'aaa'

    def test_calc_real_price_return_None(self):
        not_in_price = {'product_code': 'aaa', 'point': 1000}
        not_in_point = {'product_code': 'aaa', 'price': 1000}
        client = RakutenCrawler()
        result = client._calc_real_price(not_in_price)
        assert result == None
        result = client._calc_real_price(not_in_point)
        assert result == None

    def test_generate_next_page_query(self):
        path = os.path.join(dirname, 'product_list_page.html')
        with open(path, 'r') as f:
            response = f.read()
        client = RakutenCrawler()
        result = client._generate_next_page_query(response, {})
        assert result == {'p': '2', 'sid': '363461'}

    def test_generate_next_page_query_max_price(self):
        client = RakutenCrawler()
        response = MagicMock()
        response.url = 'https://google.com/?p=150&used=0&max=1&sid=888'
        response.text = ''
        last_product = {'price': 100}
        result = client._generate_next_page_query(response, last_product)
        assert result == {'p': '1', 'max': 100, 'sid': '888', 'used': '0'}

    def test_generate_next_page_query_faild(self):
        client = RakutenCrawler()
        response = MagicMock()
        response.url = 'https://google.com/?p=1'
        response.text = ''
        result = client._generate_next_page_query(response, {})
        assert result == None


class ScrapeDetailProductPage(unittest.TestCase):

    def test_scrape_jan_code_success(self):
        
        html_path = os.path.join(dirname, 'scrape_jan_code.html')
        response = MagicMock()
        with open(html_path, 'r') as f:
            response.text = f.read()
        response.url = 'https://item.rakuten.co.jp/superdeal/10052sp-4d151191004/'

        parsed_value = RakutenHTMLPage.scrape_product_detail_page(response)
        self.assertEqual(parsed_value.get('jan'), '4589919807796')
        self.assertEqual(parsed_value.get('price'),16000) 
        self.assertEqual(parsed_value.get('is_stocked'), True)
        self.assertEqual(parsed_value.get('product_code'), '10052sp-4d151191004')

    def test_scrape_product_detail_page_fail(self):
        response = MagicMock()
        response.text = ''
        response.url = ''
        parsed_value = RakutenHTMLPage.scrape_product_detail_page(response)
        assert parsed_value.get('jan') == None
        assert parsed_value.get('price') == None
        assert parsed_value.get('is_stocked') == False
        assert parsed_value.get('product_code') == None
        
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

    def test_parse_shop_id(self):
        path = os.path.join(dirname, 'parse_shop_id.html')
        with open(path, 'r') as f:
            response = f.read()

        shop_id = RakutenHTMLPage.parse_shop_id(response)
        assert shop_id == 384756

    def test_parse_shop_id_2(self):
        path = os.path.join(dirname, 'parse_shop_id2.html')
        with open(path, 'r') as f:
            response = f.read()

        shop_id = RakutenHTMLPage.parse_shop_id(response)
        assert shop_id == 197844

    def test_parse_next_page_url(self):
        path = os.path.join(dirname, 'product_list_page.html')
        with open(path, 'r') as f:
            response = f.read()
        result = RakutenHTMLPage.parse_next_page_url(response)
        assert result['url'] == "https://search.rakuten.co.jp/search/mall/?p=2&sid=363461"

    def test_parse_next_page_url_fail(self):
        result = RakutenHTMLPage.parse_next_page_url('')
        assert result['url'] == None
        