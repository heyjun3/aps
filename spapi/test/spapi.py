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
