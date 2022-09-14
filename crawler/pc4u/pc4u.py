import re
from datetime import datetime
import time
import urllib
import json

from bs4 import BeautifulSoup

from crawler.pc4u.models import Pc4uProduct
from crawler import utils
import settings
import log_settings
from mq import MQ


logger = log_settings.get_logger(__name__)


class CrawlerPc4u():

    def __init__(self, url: str, timestamp: datetime = datetime.now(), queue_name: str = 'mws'):
        self.url = url
        self.pc4u_product_list = []
        self.timestamp = timestamp.strftime('%Y%m%d_%H%M%S')
        self.mq = MQ(queue_name)

    def start_crawler(self):
        self.pool_product_list_page()
        self.pool_product_detail_page()

    def pool_product_list_page(self, interval_sec: int = 2):
        logger.info('action=pool_product_list_page status=run')

        while self.url is not None:
            logger.info(self.url)
            response = utils.request(self.url)
            self.pc4u_product_list.extend(Pc4uHTMLPage.scrape_product_list_page(response.text))
            self.url = Pc4uHTMLPage.scrape_next_page_url(response.text)
            time.sleep(interval_sec)

        logger.info('action=pool_product_list_page status=done')

    def pool_product_detail_page(self, interval_sec: int = 2, path: str = 'shopdetail'):
        logger.info('action=pool_product_detail_page status=run')

        for pc4u_product in self.pc4u_product_list:
            product = Pc4uProduct.get_product_code(pc4u_product.product_code)
            if product is None:
                url = urllib.parse.urljoin(settings.PC4U_ENDPOINT, f'{path}/{pc4u_product.product_code}')
                response = utils.request(url)
                parsed_value = Pc4uHTMLPage.scrape_product_detail_page(response.text)
                pc4u_product.jan = parsed_value.get('jan')
                pc4u_product.save()
                time.sleep(interval_sec)
            else:
                pc4u_product.jan = product.jan
            
            if pc4u_product.jan:
                self.publish_queue(pc4u_product.jan, pc4u_product.price)

        logger.info('action=detail_page_crawling status=done')

    def publish_queue(self, jan: str, price: int):
        logger.info('action=publish_queue status=run')

        self.mq.publish(json.dumps({
            'filename': f'pc4u_{self.timestamp}',
            'jan': jan,
            'cost': price
        }))

        logger.info('action=publish_queue status=done')


class Pc4uHTMLPage(object): 

    @staticmethod
    def scrape_product_list_page(response: str) -> list[Pc4uProduct]:
        logger.info('action=scrape_product_list_page status=run')

        pc4u_product_list = []
        soup = BeautifulSoup(response, 'lxml')
        products = soup.select('.innerBox')

        for product in products:
            title = product.select_one('.name a')
            cost = product.select_one('.price')
            if not title or not cost:
                continue
            try:
                product_code = re.search('[0-9]{12}', title.attrs.get('href')).group()
            except AttributeError as ex:
                logger.info(f'product_code is None error={ex}')
                continue

            title = title.text.strip()
            cost = int(''.join(re.findall('[\\d+]', cost.text)))

            pc4u_product = Pc4uProduct(name=title, price=cost, product_code=product_code)
            pc4u_product_list.append(pc4u_product)
        
        logger.info('action=scrape_product_list_page status=done')
        return pc4u_product_list

    @staticmethod
    def scrape_product_detail_page(response: str) -> dict:
        logger.info('action=scrape_product_detail_page status=run')

        soup = BeautifulSoup(response, 'lxml')
        detail_text = soup.select_one('.detailTxt')
        try:
            jan = re.search('[0-9]{13}', str(detail_text)).group()
        except AttributeError as ex:
            logger.info(f'jan is None error={ex}')
            jan = None

        try:
            price = soup.select_one('#M_price1').attrs.get('value')
            price = int(''.join(re.findall('[0-9]', price)))
        except (AttributeError, TypeError) as ex:
            logger.info({'message': 'price is None', 'error': ex})
            price = None

        is_stocked = soup.select_one('.cartBtn')

        logger.info('action=scrape_product_detail_page status=done')
        return {'jan': jan, 'price': price, 'is_stocked': bool(is_stocked)}

    @staticmethod
    def scrape_next_page_url(response: str) -> str|None:
        logger.info('action=scrape_next_page_url status=run')

        soup = BeautifulSoup(response, 'lxml')
        try:
            next_url = soup.select('.M_pager li')[-1].a.attrs.get('href')
        except AttributeError as ex:
            logger.info(ex)
            return None

        sold_out_flag = soup.select_one('.innerBox .btnWrap img')
        if sold_out_flag and sold_out_flag.attrs['alt'] == '品切れ':
            logger.info('next page url is None')
            return None

        next_url = urllib.parse.urljoin(settings.PC4U_ENDPOINT, next_url)
        return next_url

def main():
    outlet_url = 'https://www.pc4u.co.jp/shopbrand/outlet/'
    sale_url = 'https://www.pc4u.co.jp/shop/shopbrand.html?search=&prize1=sale'
    new_url = 'https://www.pc4u.co.jp/shop/shopbrand.html?search=new'
    url = "https://www.pc4u.co.jp/shop/shopbrand.html?page=1&search="
    urls = [outlet_url, sale_url, new_url, url]

    timestamp = datetime.now()
    for u in urls:
        client = CrawlerPc4u(url=u, timestamp=timestamp)
        client.start_crawler()
