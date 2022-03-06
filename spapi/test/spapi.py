import unittest
import os
import json

from spapi.spapi import SPAPIJsonParser


dirname = os.path.join(os.path.dirname(__file__), 'test_json')

class ParseGetCompetitivePricing(unittest.TestCase):

    def test_success_parse(self):
        path = os.path.join(dirname, 'parse_get_competitive_pricing.json')
        with open(path, 'r') as f:
            response = f.read()

        products = SPAPIJsonParser.parse_get_competitive_pricing(json.loads(response))
        
        self.assertEqual(products[0]['asin'], "B08HMT3LRN")
        self.assertEqual(products[0]['price'], 3267)
        self.assertEqual(products[0]['ranking'], 74218)
        self.assertEqual(products[-1]['asin'], "B07HG6F6K2")
        self.assertEqual(products[-1]['price'], -1)
        self.assertEqual(products[1]['ranking'], -1)


class ParseGetItemOffers(unittest.TestCase):

    def test_success_parse(self):
        path = os.path.join(dirname, 'get_item_offers.json')
        with open(path, 'r') as f:
            response = f.read()

        product = SPAPIJsonParser.parse_get_item_offers(json.loads(response))

        self.assertEqual(product['asin'], 'B07HG6F6K2')
        self.assertEqual(product['price'], 2800)
        self.assertEqual(product['ranking'], 15)


class ParseListCatalogItems(unittest.TestCase):

    def test_success_parse(self):
        path = os.path.join(dirname, 'list_catalog_items.json')
        with open(path, 'r') as f:
            response = f.read()

        products = SPAPIJsonParser.parse_list_catalog_items(json.loads(response))
        self.assertEqual(products[0]['asin'], "B00131HAQC")
        self.assertEqual(products[0]['quantity'], 1)
        self.assertEqual(products[0]['title'], "アロン化成 安寿 ポータブルトイレ用防臭剤22")
        self.assertEqual(products[0]['price'], 1760)

        self.assertEqual(products[-1]['asin'], "B08R16SCR5")
        self.assertEqual(products[-1]['quantity'], 8)
        self.assertEqual(products[-1]['title'], "【まとめ買い】アロン化成 安寿 ポータブルトイレ用防臭剤22×8個")
        self.assertEqual(products[-1]['price'], None)

    def test_response_is_none(self):
        path = os.path.join(dirname, 'list_catalog_items_none.json')
        with open(path, 'r') as f:
            response = f.read()

        products = SPAPIJsonParser.parse_list_catalog_items(json.loads(response))
        self.assertFalse(products)

