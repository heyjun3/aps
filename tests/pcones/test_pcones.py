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

    def test_scrape_detail_page_1(self):
        path = os.path.join(dirname, 'scrape_product_detail_page1.html')
        with open(path, 'r') as f:
            response = f.read()
        
        parsed_value = PconesHTMLPage.scrape_product_detail_page(response)
        assert parsed_value.get('jan') == '4589967503374'
        assert parsed_value.get('price') == 13970
        assert parsed_value.get('is_stocked') == False

    def test_scrape_detail_page_2(self):
        path = os.path.join(dirname, 'scrape_product_detail_page2.html')
        with open(path, 'r') as f:
            response = f.read()

        parsed_value = PconesHTMLPage.scrape_product_detail_page(response)
        assert parsed_value.get('jan') == None
        assert parsed_value.get('price') == 3499
        assert parsed_value.get('is_stocked') == True
