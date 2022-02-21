from turtle import ht
import unittest
import unittest.mock
import os

from crawler.pc4u.pc4u import Pc4uHTMLPage


dirname = os.path.join(os.path.dirname(__file__), 'test_html')



class ScrapeProductListPage(unittest.TestCase):

    def test_scrape_product_list_page(self):
        response = unittest.mock.MagicMock()
        html_page = os.path.join(dirname, 'scrape_product_list_page.html')
        with open(html_page, 'r') as f:
            response.text = f.read()

        products = Pc4uHTMLPage.scrape_product_list_page(response.text)
        self.assertEqual(products[0].name, '【アウトレット特価・新品】Cooler Master MasterCase H500M E-ATX ミドルタワー型PCケース｜MCM-H500M-IHNN-S00')
        self.assertEqual(products[0].price, 21340)
        self.assertEqual(products[0].product_code, '000000078117')

        self.assertEqual(products[-1].name, '【アウトレット特価・新品】Arashi Vision GO2 レンズ保護フィルター｜CING2CB/B')
        self.assertEqual(products[-1].price, 1140)
        self.assertEqual(products[-1].product_code, '000000073973')


class ScrapeProductDetailPage(unittest.TestCase):

    def test_scrape_detail_product_page(self):
        response = unittest.mock.MagicMock()
        html_page = os.path.join(dirname, 'scrape_detail_product_page.html')
        with open(html_page, 'r') as f:
            response.text = f.read()

        jan = Pc4uHTMLPage.scrape_product_detail_page(response.text)
        
        self.assertEqual(jan, '4537694274043')


class ScrapeNextPageUrl(unittest.TestCase):

    def test_scrape_next_page_url(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'scrape_next_page_url.html')
        with open(html_path, 'r') as f:
            response.text = f.read()

        url = Pc4uHTMLPage.scrape_next_page_url(response.text)
        self.assertEqual(url, 'https://www.pc4u.co.jp/shopbrand/outlet/page2/order/')

    def test_soled_out(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'soled_out.html')
        with open(html_path, 'r') as f:
            response.text = f.read()

        url = Pc4uHTMLPage.scrape_next_page_url(response.text)
        self.assertIsNone(url)

    def test_last_page(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'last_page.html')
        with open(html_path, 'r') as f:
            response.text = f.read()

        url = Pc4uHTMLPage.scrape_next_page_url(response.text)
        self.assertIsNone(url)
