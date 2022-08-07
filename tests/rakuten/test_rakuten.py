import unittest
import os

from crawler.rakuten.rakuten import RakutenHTMLPage


dirname = os.path.join(os.path.dirname(__file__), 'test_html')


class ScrapeDetailProductPage(unittest.TestCase):

    def test_scrape_jan_code(self):
        
        html_path = os.path.join(dirname, 'scrape_jan_code.html')
        with open(html_path, 'r') as f:
            response = f.read()

        jan = RakutenHTMLPage.scrape_product_detail_page(response)
        self.assertEqual(jan, '4573201242433')
        