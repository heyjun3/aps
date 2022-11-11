import os
import json
import datetime

import pytest

from crawler.paypaymall.paypaymoll import PayPayMollHTMLParser
from crawler.paypaymall.paypaymoll import YahooShopApiParser
from crawler.paypaymall import paypaymoll

dirname = os.path.join(os.path.dirname(__file__), 'test_html')


class TestPayPayMoll(object):

    def test_product_detail_page(self):
        path = os.path.join(dirname, 'product_detail_page.html')
        with open(path, 'r') as f:
            response = f.read()

        parsed_value = PayPayMollHTMLParser.product_detail_page_parser(response)
        assert parsed_value.get('jan') == '4957054501273'
        assert parsed_value.get('price') == 38500
        assert parsed_value.get('is_stocked') == False

class TestYahooShopCrawler(object):

    def test_calc_real_price(self):
        result = paypaymoll.ItemSearchResult('test', 1000, '4444', 'test', 100, 'id', 'url')
        client = paypaymoll.YahooShopCrawler()
        value = client._calc_real_price(result)
        assert value == paypaymoll.ItemSearchResult('test', 900, '4444', 'test', 100, 'id', 'url')

    def test_calc_real_price_fail(self):
        result = paypaymoll.ItemSearchResult(None, None)
        client = paypaymoll.YahooShopCrawler()
        value = client._calc_real_price(result)
        assert value == None

        value = client._calc_real_price(None)
        assert value == None

    def test_generate_publish_message(self):
        result = paypaymoll.ItemSearchResult('test', 1000, '4444', 'test', 100, 'id', 'url')
        client = paypaymoll.YahooShopCrawler()
        timestamp = datetime.datetime(2022, 11, 11, 0, 49, 1)
        value = client._generate_publish_message(result, timestamp)
        assert value == '{"jan": "4444", "cost": 1000, "url": "url", "filename": "paypay_20221111_004901"}'

    def test_generate_publish_message(self):
        client = paypaymoll.YahooShopCrawler()
        value_1 = paypaymoll.ItemSearchResult('test', None, '4444', 'test', 100, 'id', 'url')
        value_2 = paypaymoll.ItemSearchResult('test', 1000, None, 'test', 100, 'id', 'url')
        value_3 = paypaymoll.ItemSearchResult('test', 1000, '4444', 'test', 100, 'id', None)
        result_1 = client._generate_publish_message(value_1, datetime.datetime.now())
        result_2 = client._generate_publish_message(value_2, datetime.datetime.now())
        result_3 = client._generate_publish_message(value_3, datetime.datetime.now())
        result_4 = client._generate_publish_message(None, datetime.datetime.now())
        assert result_1 == None
        assert result_2 == None
        assert result_3 == None
        assert result_4 == None


class TestYahooShopApiParser(object):

    def test_parse_item_search_v3_success(self):
        path = os.path.join(dirname, 'item_search_v3.json')
        with open(path, 'r') as f:
            res = f.read()

        value = YahooShopApiParser.parse_item_search_v3(json.loads(res))
        assert len(value) == 100
        assert value[0].product_id == "ksdenki_4549980503980"
        assert value[0].price == 100000
        assert value[0].jan == "4549980503980"
        assert value[0].name == "Panasonic\uff08\u30d1\u30ca\u30bd\u30cb\u30c3\u30af\uff09 \u6b21\u4e9c\u5869\u7d20\u9178\u7a7a\u9593\u9664\u83cc\u8131\u81ed\u6a5f\u3000\u30b8\u30a2\u30a4\u30fc\u30ce F-MV2300-WZ"
        assert value[0].point == 1000
        assert value[0].shop_id == "ksdenki"
        assert value[0].url == "https://store.shopping.yahoo.co.jp/ksdenki/4549980503980.html"

        assert value[-1].product_id == "ksdenki_4580053310029"
        assert value[-1].price == 89100
        assert value[-1].jan == "4580053310029"
        assert value[-1].name == "PowerVision\uff08\u30d1\u30ef\u30fc\u30d3\u30b8\u30e7\u30f3\uff09 \u6c34\u4e2d\u30c9\u30ed\u30fc\u30f3\u3000PowerRay PRE10"
        assert value[-1].point == 5346
        assert value[-1].shop_id == "ksdenki"
        assert value[-1].url == "https://store.shopping.yahoo.co.jp/ksdenki/4580053310029.html"

    def test_parse_item_search_v3_fail(self):
        with pytest.raises(TypeError) as ex:
            value = YahooShopApiParser.parse_item_search_v3({})
            
        assert str(ex.value) == str(TypeError("'NoneType' object is not iterable"))
