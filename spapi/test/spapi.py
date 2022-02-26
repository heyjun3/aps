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