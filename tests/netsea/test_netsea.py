import unittest
import unittest.mock
import os

from crawler.netsea.netsea import NetseaHTMLPage


dirname = os.path.join(os.path.dirname(__file__), 'test_html')


class ScrapeNextPageUrl(unittest.TestCase):

    def test_scrape_next_page_url(self):
        response = unittest.mock.MagicMock()
        response.url = 'https://www.netsea.jp/search?category_id=1&ex_so=N&sort=new&searched=Y&page=166'
        html_path = os.path.join(dirname, 'netsea_scrape_next_page.html')
        with open(html_path, 'r') as f:
            response.text = f.read()

        url = NetseaHTMLPage.scrape_next_page_url(response.text, response.url, is_new_product_search=True)
        self.assertIsNone(url)

    def test_scrape_next_page_url_for_favorite_list_first_page(self):
        response = unittest.mock.MagicMock()
        response.url = 'https://www.netsea.jp/bookmark'
        html_path = os.path.join(dirname, 'first_favorite_product_list_page.html')
        with open(html_path, 'r') as f:
            response.text = f.read()

        next_page_url = 'https://www.netsea.jp/bookmark?page=2'        
        url = NetseaHTMLPage.scrape_next_page_url(response.text, response.url)
        self.assertEqual(next_page_url, url)

    def test_scrape_next_page_url_for_favorite_list_last_page(self):
        response = unittest.mock.MagicMock()
        response.url = 'https://www.netsea.jp/bookmark?page=23'
        html_path = os.path.join(dirname, 'last_favorite_product_list_page.html')
        with open(html_path, 'r') as f:
            response.text = f.read()

        url = NetseaHTMLPage.scrape_next_page_url(response.text, response.url)
        self.assertIsNone(url)

    def test_scrape_change_price(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'scrape_change_price.html')
        response.url = 'https://www.netsea.jp/search/?facet_price_to=4508&disc_flg=Y&ex_so=Y&sort=PD&searched=Y&page=166'
        with open(html_path, 'r') as f:
            response.text = f.read()

        url = NetseaHTMLPage.scrape_next_page_url(response.text, response.url)
        self.assertEqual(url, 'https://www.netsea.jp/search/?facet_price_to=4098&disc_flg=Y&ex_so=Y&sort=PD&searched=Y&page=1')


class ScrapeProductListPage(unittest.TestCase):

    def test_scrape_discount_price(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'scrape_discount_price.html')
        with open(html_path, 'r') as f:
            response.text = f.read()
        
        products = NetseaHTMLPage.scrape_product_list_page(response.text)
        self.assertEqual(products[0].price, 6049999998)
        self.assertEqual(products[0].name, '非売品　注文しないてください' )
        self.assertEqual(products[0].shop_code, '750972')
        self.assertEqual(products[0].product_code, 'FMP001')
        self.assertEqual(products[0].url, 'https://www.netsea.jp/shop/750972/FMP001?cx_search=PD&fw=ChQKEgoQCKWp9cWVi_YCFe2GBgoYHA')

        self.assertEqual(products[-1].price, 294030)
        self.assertEqual(products[-1].name, '☆原石一点物☆【原石】アメジスト カペーラ (3A) (ウルグアイ産) (台付) (22.5kg) No.39' )
        self.assertEqual(products[-1].shop_code, '172479')
        self.assertEqual(products[-1].product_code, 'FE09Jmi064-39')
        self.assertEqual(products[-1].url,'https://www.netsea.jp/shop/172479/FE09Jmi064-39?cx_search=PD&fw=ChYKEgoQCKWp9cWVi_YCFe2GBgoYHBA7')

    def test_scrape_not_discount_price(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'scrape_not_discount_price.html')
        with open(html_path, 'r') as f:
            response.text = f.read()
        
        products = NetseaHTMLPage.scrape_product_list_page(response.text)
        self.assertEqual(products[0].price, 33611)
        self.assertEqual(products[0].name, '[直送品]リッチェル コアラクーン 両対面式A型ベビーカー 1カ月頃から36カ月頃 ネイビーブルー' )
        self.assertEqual(products[0].shop_code, '5984')
        self.assertEqual(products[0].product_code, '4973655937907')
        self.assertEqual(products[0].url, 'https://www.netsea.jp/shop/5984/4973655937907?cx_search=PD&fw=ChQKEgoQCL37j7ugi_YCFYuiBgoYHA')

        self.assertEqual(products[-1].price, 6600)
        self.assertEqual(products[-1].name, 'カプセル粉づめくん　本体　０号用' )
        self.assertEqual(products[-1].shop_code, '5984')
        self.assertEqual(products[-1].product_code, '4905712000521')
        self.assertEqual(products[-1].url, 'https://www.netsea.jp/shop/5984/4905712000521?cx_search=PD&fw=ChYKEgoQCL37j7ugi_YCFYuiBgoYHBA7')


class ScrapeProductDetailPage(unittest.TestCase):

    def test_scrape_jan_code(self):
        response = unittest.mock.MagicMock()
        html_path = os.path.join(dirname, 'scrape_jan_code.html')
        with open(html_path, 'r') as f:
            response.text = f.read()

        parsed_value = NetseaHTMLPage.scrape_product_detail_page(response.text)
        self.assertEqual(parsed_value.get('jan'), '4962644942725')

class TestScrapeFavoriteListPage(object):

    def test_scrape_favorite_list_page(self):
        html_path = os.path.join(dirname, "scrape_favorite_list_page.html") 
        with open(html_path, "r") as f:
            file = f.read()
        
        parsed = NetseaHTMLPage.scrape_favorite_list_page(file)

        assert len(parsed) == 40
        assert parsed[0].name == "NEW日本シグマックス メディエイド しっかりガード 腰 スタンダードプラス L"
        assert parsed[0].price == 2300
        assert parsed[0].shop_code == "84918"
        assert parsed[0].product_code == "28234210"
        assert parsed[0].url == "https://www.netsea.jp/shop/84918/28234210"
        assert parsed[-1].name == "韓国コスメ エイプリルスキン（APRILSKIN）リアルカレン ピールオフパック　100ｇ"
        assert parsed[-1].price == 1019
        assert parsed[-1].shop_code == "8070"
        assert parsed[-1].product_code == "2022020703"
        assert parsed[-1].url == "https://www.netsea.jp/shop/8070/2022020703"
