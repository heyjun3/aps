import os
import unittest
import unittest.mock

from crawler.super import super
from crawler.super.super import SuperHTMLPage

dirname = os.path.join(os.path.dirname(__file__), 'test_html')


class ScrapeProductListPage(unittest.TestCase):

    def test_scrape_product_list_page(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'scrape_product_list_page.html')
        url = 'https://www.superdelivery.com/p/do/dpsl/21321/?so=newly'
        with open(html_path, 'r') as f:
            response.text = f.read()

        products = SuperHTMLPage.scrape_product_list_page(response.text, url)

        self.assertEqual(products[0].name, '【特価】ベジート・ボヌール IH対応二食鍋30cm（ARB-2212）')
        self.assertEqual(products[0].shop_code, '21321')
        self.assertEqual(products[0].product_code, '10032806')
        self.assertEqual(products[0].price, 2783)
        
        self.assertEqual(products[-1].name, 'フロウト スポンジラック 置き型（RG-0461）')
        self.assertEqual(products[-1].shop_code, '21321')
        self.assertEqual(products[-1].product_code, '9925316')
        self.assertEqual(products[-1].price, 792)


class ScrapeProductDetailPage(unittest.TestCase):

    def test_scrape_product_detail_page(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'scrape_product_detail_page.html')
        with open(html_path, 'r') as f:
            response.text = f.read()
        products = SuperHTMLPage.scrape_product_detail_page(response.text)

        self.assertEqual(products[0].jan, '4903779122125')
        self.assertEqual(products[0].price, 2783)
        self.assertEqual(products[0].set_number, 1)
        self.assertEqual(products[0].shop_code, '21321')
        self.assertEqual(products[0].product_code, '10032806')

        self.assertEqual(products[-1].jan, '4903779122125')
        self.assertEqual(products[-1].price, 2783)
        self.assertEqual(products[-1].set_number, 2)
        self.assertEqual(products[-1].shop_code, '21321')
        self.assertEqual(products[-1].product_code, '10032806')


class ScrapeShopListPage(unittest.TestCase):

    def test_scrape_shop_list_page(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'scrape_shop_list_page.html')
        with open(html_path, 'r') as f:
            response.text = f.read()

        shops = SuperHTMLPage.scrape_shop_list_page(response.text)

        self.assertEqual(shops[0].name, 'ProjectID')
        self.assertEqual(shops[0].shop_id, '1003514')

        self.assertEqual(shops[-1].name, '夢家ハチマン')
        self.assertEqual(shops[-1].shop_id, '1002934')


class ScrapeNextPageUrl(unittest.TestCase):

    def test_scrape_shop_list_page(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'scrape_shop_list_page.html')
        with open(html_path, 'r') as f:
            response.text = f.read()
        url = SuperHTMLPage.scrape_next_page_url(response.text)
        self.assertEqual(url, 'https://www.superdelivery.com/p/do/psl/all/2/?so=newdealer')

    def test_scrape_product_list_page(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'scrape_product_list_page.html')
        with open(html_path, 'r') as f:
            response.text = f.read()
        url = SuperHTMLPage.scrape_next_page_url(response.text)
        self.assertEqual(url, 'https://www.superdelivery.com/p/do/dpsl/21321/all/2/?so=newly')
