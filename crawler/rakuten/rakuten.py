""""
This file is Rakuten api main file

"""
import re
import time
import datetime
import os

import requests
import openpyxl
from bs4 import BeautifulSoup

from crawler.rakuten.models import RakutenProduct
import settings
import log_settings

SHOP_CODES = ['ksdenki', 'dj', 'e-zoa', 'reckb', 'jtus', 'ioplaza', 'ikebe']

logger = log_settings.get_logger(__name__)


def jan_selector(item: dict) -> list:
    """Search jan from itemInfo"""
    item_url = item['Item']['itemUrl']
    jan = re.findall('[0-9]{13}', item_url)

    if jan is None:
        item_caption = item['Item']['itemCaption']
        jan = re.findall('[0-9]{13}', item_caption)

    if jan:
        jan = jan[0]
    return jan


def product_page_parser(response: str):
    """product_page parse return jan_code
    if jan_code is None return None object"""
    logger.info('action=product_page_parser status=run')

    soup = BeautifulSoup(response, 'lxml')
    jan = soup.select_one('.item_number')
    if jan is None:
        logger.info("product page hasn't product_code")
        return None
    jan = re.findall('[0-9]{13}', jan.text)
    if not jan:
        logger.info("product code isn't jan code")
        return None
    logger.info(jan[0])
    return jan[0]


class Rakuten:
    """Rakuten api Class"""
    def __init__(self, shop_code: str):
        self.shop_code = shop_code
        self.products = []
        self.params = {
            'applicationId': settings.RAKUTEN_APP_ID,
            'shopCode': self.shop_code,
            'page': 1,
            'sort': '-itemPrice',
            'maxPrice': None,
        }

    def main(self):
        """running Rakuten api method"""
        logger.info('action=main status=run')
        while True:
            flag_info = self.rakuten_api()
            if flag_info['item_count'] < 30:
                break
            if self.params['page'] == 100:
                self.params['page'] = 1
                if self.params['maxPrice'] == flag_info['last_price']:
                    flag_info['last_price'] -= 100
                self.params['maxPrice'] = flag_info['last_price']

            self.params['page'] += 1

    def rakuten_api(self) -> dict:
        """rakuten api method"""
        logger.info('action=rakuten_api status=run')
        logger.info(f'action=rakuten_api page={self.params["page"]} maxPrice={self.params["maxPrice"]}')
        response = self.rakuten_request()
        logger.info(f'action=rakuten_api status_code={response.status_code}')
        response = response.json()
        time.sleep(2)

        last_price = self.products_info_selector(response)
        flag_info = {'item_count': len(response['Items']), 'last_price': last_price}

        return flag_info

    def rakuten_request(self):
        """rakuten api request method"""
        for _ in range(60):
            try:
                response = requests.get(settings.REQUEST_URL, timeout=30.0, params=self.params)
                if not response.status_code == 200 or response is None:
                    raise Exception
                return response
            except Exception as e:
                logger.error(f'action=request error={e}')
                time.sleep(30)

    def products_info_selector(self, response: dict):
        """requests response search jan and price.
            self.info_list add jan and price"""
        last_price = None

        for item in response['Items']:
            price = item['Item']['itemPrice']
            point_rate = item['Item']['pointRate']
            last_price = price
            calc_price = int(int(price) * (91 - int(point_rate)) / 100)
            item_name = item['Item']['itemName']
            url = item['Item']['itemUrl']
            jan = jan_selector(item)
            product = RakutenProduct.create(name=item_name, jan=jan, price=calc_price,
                                            shop_code=self.shop_code, url=url)
            self.products.append(product)
        return last_price

    def export_to_excel(self):
        """self.jan_price_list export to _excel"""
        logger.info('action=export_to_excel status=run')
        if not self.products:
            logger.warning('action=export_to_excel error=jan_list is None')
            return

        workbook = openpyxl.Workbook()
        sheet = workbook['Sheet']
        sheet.append(['JAN', 'Cost'])
        for product in self.products:
            if product.jan and product.price:
                sheet.append([product.jan, product.price])
            else:
                logger.error(product.value)

        timestamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')
        save_path = os.path.join(settings.SCRAPE_SAVE_PATH, f'rakuten{timestamp}.xlsx')
        workbook.save(save_path)
        workbook.close()


def main(shop_code: str):
    logger.info('action=rakuten_main status=run')
    rakuten = Rakuten(shop_code=shop_code)
    rakuten.main()
    logger.info(rakuten.products)
    for product in rakuten.products:
        product.get_jan_code()
        logger.info(product.value)
    rakuten.export_to_excel()


def schedule():
    for shop_code in SHOP_CODES:
        main(shop_code)
