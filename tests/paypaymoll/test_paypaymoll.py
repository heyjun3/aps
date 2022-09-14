import os

from crawler.paypaymall.paypaymoll import PayPayMollHTMLParser

dirname = os.path.join(os.path.dirname(__file__), 'test_html')


class TestPayPayMoll(object):

    def test_product_detail_page(self):
        path = os.path.join(dirname, 'product_detail_page.html')
        with open(path, 'r') as f:
            response = f.read()

        parsed_value = PayPayMollHTMLParser.product_detail_page_parser(response)
        assert parsed_value.get('jan') == '4957054501273'
        assert parsed_value.get('price') == 38500
        assert parsed_value.get('is_stocked') == False
