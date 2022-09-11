from ast import parse
import unittest
import unittest.mock
import os

from crawler.buffalo.buffalo import BuffaloHTMLPage


dirname = os.path.join(os.path.dirname(__file__), 'test_html')

class ScrapeProductListPage(unittest.TestCase):

    def test_scrape_product_list_page(self):
        response = unittest.mock.MagicMock
        html_path = os.path.join(dirname, 'scrape_product_list_page.html')
        with open(html_path, 'r') as f:
            response.text = f.read()

        products = BuffaloHTMLPage.scrape_product_list_page(response.text)
        
        self.assertEqual(products[0].name, '《アウトレット・未使用》HD-QHA32U3/R5(保証有り)')
        self.assertEqual(products[0].price, 177800)
        self.assertEqual(products[0].product_code, '25444')

        self.assertEqual(products[-1].name, '《アウトレット・未使用》CF-BRIDGE(保証有り)')
        self.assertEqual(products[-1].price, 620)
        self.assertEqual(products[-1].product_code, '23327')


class ScrapeProductDetailPage(unittest.TestCase):

    def test_scrape_product_detail_page(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'scrape_product_detail_page.html')
        with open(html_path, 'r') as f:
            response.text = f.read()

        parsed_value = BuffaloHTMLPage.scrape_product_detail_page(response.text)
        self.assertEqual(parsed_value.get('jan'), '4988755224604')
        self.assertEqual(parsed_value.get('price'), 620)
        self.assertEqual(parsed_value.get('is_stocked'), True)

    # def test_product_jan_code_is_none(self):
    #     response = unittest.mock.MagicMock()
    #     html_path = os.path.join(dirname, 'product_jan_code_is_none.html')
    #     with open(html_path, 'r') as f:
    #         response.text = f.read()

    #     jan = BuffaloHTMLPage.scrape_product_detail_page(response.text)
    #     self.assertIsNone(jan)