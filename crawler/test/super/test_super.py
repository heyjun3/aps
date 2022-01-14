import os
import unittest
import unittest.mock

from crawler.super import super


class TestSuper(unittest.TestCase):

    def test_detail_page_selector(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(os.path.dirname(__file__), 'test_super_product_detail_page.html')
        with open(html_path, 'r') as f:
            response.text = f.read()

        products = super.detail_page_selector(response)

        self.assertEqual('9037265', products[0].product_code)
        self.assertEqual('VSTN-2000B', products[0].product_detail_code)
        self.assertEqual('184955', products[0].shop_code)
        self.assertEqual(3960, products[0].price)
        self.assertEqual('4977642021907', products[0].jan)
