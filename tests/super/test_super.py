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
        self.assertEqual(products[0].url, 'https://www.superdelivery.com/p/r/pd_p/10032806/')
        
        self.assertEqual(products[-1].name, 'フロウト スポンジラック 置き型（RG-0461）')
        self.assertEqual(products[-1].shop_code, '21321')
        self.assertEqual(products[-1].product_code, '9925316')
        self.assertEqual(products[-1].price, 792)
        self.assertEqual(products[-1].url, 'https://www.superdelivery.com/p/r/pd_p/9925316/')


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

    def test_scrape_favorite_products_page(self):
        path = os.path.join(dirname, "favorite_product_page.html")
        with open(path, "r") as f:
            res = f.read()

        url = SuperHTMLPage.scrape_next_page_url(res)

        assert url == "https://www.superdelivery.com/p/wishlist/search.do?o=ad&p=2&cc=all&dc=0"

    def test_scrape_last_favorite_products_page(self):
        path = os.path.join(dirname, "favorite_last_page.html")
        with open(path, "r") as f:
            res = f.read()
        
        url = SuperHTMLPage.scrape_next_page_url(res)

        assert url == None

class TestScrapeFavoriteProductListPage(object):

    def test_scrape_products(self):
        path = os.path.join(dirname, "favorite_product_page.html")
        with open(path, "r") as f:
            res = f.read()
        
        products = SuperHTMLPage.scrape_favorite_product_list_page(res)

        assert len(products) == 48
        assert products[0].name == "80%OFF【スクエアエニックス】ドラゴンクエスト 文具屋 冒険ダイアリー2023"
        assert products[0].product_code == "11049757"
        assert products[0].price == 655
        assert products[0].url == "https://www.superdelivery.com/p/r/pd_p/11049757/"
        assert products[0].shop_code == None
        assert products[-1].name == "ビクトリノックス ( VICTORINOX )　1.3713　ハントマン　91mm　レッド　箱入"
        assert products[-1].product_code == "10182027"
        assert products[-1].price == 4287
        assert products[-1].url == "https://www.superdelivery.com/p/r/pd_p/10182027/"
        assert products[-1].shop_code == None
