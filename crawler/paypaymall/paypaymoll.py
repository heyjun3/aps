import re

from bs4 import BeautifulSoup

import log_settings


logger = log_settings.get_logger(__name__)


class PayPayMollHTMLParser(object):

    @staticmethod
    def product_detail_page_parser(response: str) -> dict:

        soup = BeautifulSoup(response, 'lxml')
        try:
            item_details = soup.select('.ItemDetails_list')
            item_detail = list(filter(None, map(lambda x: ''.join(re.findall('[0-9]', x.text)), item_details)))
            jan = item_detail.pop() if item_detail else None
        except (AttributeError) as ex:
            logger.info({'message': 'jan is None', 'error': ex})
            jan = None

        price = soup.select_one('.ItemPrice_price')
        if price:
            price = int(''.join(re.findall('[0-9]', price.text)))

        is_stocked = soup.select_one('#CartButtonUltLog').attrs
        if 'disabled' in is_stocked:
            is_stocked = False

        return {'jan': jan, 'price': price, 'is_stocked': bool(is_stocked)}
