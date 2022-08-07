import unittest
import os

from crawler.pcones.pcones import PconesHTMLPage


dirname = os.path.join(os.path.dirname(__file__), 'test_html')


class ScrapeProductListPage(unittest.TestCase):

    def test_scrape_success(self):
        path = os.path.join(dirname, 'scrape_product_list_page.html')
        with open(path, 'r') as f:
            response = f.read()

        products = PconesHTMLPage.scrape_product_list_page(response)

        self.assertEqual(products[0]['jan'], '4545708003824')
        self.assertEqual(products[0]['cost'], 15800)

        self.assertEqual(products[-1]['jan'], '4550161189572')
        self.assertEqual(products[-1]['cost'], 6080)
