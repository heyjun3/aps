import re
import datetime
import time
import urllib
import os

from requests import Response
from bs4 import BeautifulSoup
import pandas as pd

from crawler.pc4u.models import Pc4u
from crawler import utils
import settings
import log_settings


logger = log_settings.get_logger(__name__)

PC4U_SHOPCODE = 'https://www.pc4u.co.jp'


class CrawlerPc4u():

    filename = 'pc4u'

    def __init__(self, url):
        self.url = url
        self.products = []

    def list_page_crawling(self, interval_sec: int = 2):
            logger.info('action=list_page_crawling status=run')

            while True:
                response = utils.request(self.url)
                self.list_page_scraping(response)
                time.sleep(interval_sec)
                if not self.url:
                    logger.info('action=list_page_crawling status=done')
                    break
                self.scrape_next_page_url(response)

    def list_page_scraping(self, response: Response):
        logger.info('action=list_page_product_scraping status=run')

        soup = BeautifulSoup(response.text, 'lxml')
        products = soup.select('.innerBox')

        for product in products:
            title = product.select_one('.name a')
            cost = product.select_one('.price')
            if not title or not cost:
                continue
            product_code = re.search('[0-9]{12}', title.attrs['href'])
            title = title.text.strip()
            cost = int(''.join(re.findall('[\\d+]', cost.text)))
            if product_code:
                product_code = product_code[0]
            sold_out_flag = product.select_one('.btnWrap img')
            if sold_out_flag and sold_out_flag.attrs['alt'] == '品切れ':
                self.url = None
                return
            url = urllib.parse.urljoin(PC4U_SHOPCODE, f'shopdetail/{product_code}')
            pc4u = Pc4u.create(name=title, price=cost, shop_code=PC4U_SHOPCODE,
                               url=url, product_code=product_code)
            self.products.append(pc4u)
            logger.debug(pc4u.value)
        logger.info('action=list_page_product_scraping status=done')

    def scrape_next_page_url(self, response: Response):
        logger.info('action=next_url_scraping status=run')

        soup = BeautifulSoup(response.text, 'lxml')
        next_url = soup.select('.M_pager li')
        if not next_url[-1].attrs['class'][0] == 'next':
            self.url = None
            return
        next_url = next_url[-1].select_one('a').attrs['href']
        next_url = urllib.parse.urljoin(PC4U_SHOPCODE, next_url)
        logger.debug(next_url)
        self.url = next_url

    @classmethod
    def scrape_jan_code_from_detail_product_page(cls, response: Response):
        logger.info('action=scrape_jan_code_from_detail_product_page')
        soup = BeautifulSoup(response.text, 'lxml')
        for p_tag in soup.select('p'):
            if 'JAN' in p_tag.text.strip():
                jan = re.search('[0-9]{13}', p_tag.text)
                if jan:
                    logger.debug(jan[0])
                    return jan[0]

        for td_tag in soup.select('td'):
            jan = re.search('[0-9]{13}', td_tag.text)
            if jan:
                logger.debug(jan[0])
                return jan[0]

    def detail_page_crawling(self, interval_sec: int = 2):
        logger.info('action=detail_page_crawling status=run')

        for product in self.products:
            db_response = Pc4u.get_product_code(product.product_code)
            if db_response is None:
                response = utils.request(product.url)
                jan = CrawlerPc4u.scrape_jan_code_from_detail_product_page(response)
                product.jan = jan
                product.save()
                time.sleep(interval_sec)
                continue
            elif not db_response.price == product.price:
                Pc4u.price_update('product_code', product.product_code, product.price)
            product.jan = db_response.jan

        logger.info('action=detail_page_crawling status=done')

    def save_excel_file(self, save_path=settings.SCRAPE_SAVE_PATH):
        logger.info('action=save_excel_file status=run')

        timestamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')
        jan_cost_list = [[product.jan, product.price] for product in self.products if product.jan]
        df = pd.DataFrame(data=jan_cost_list, columns=['JAN', 'Cost']).astype({'JAN': str, 'Cost': int}).drop_duplicates()
        if not df.empty:
            df.to_excel(os.path.join(save_path, f'{self.filename}{timestamp}.xlsx'), index=False)

        logger.info('action=save_excel_file status=done')


def schedule_pc4u_task_everyday():
    outlet_url = 'https://www.pc4u.co.jp/shopbrand/outlet/'
    sale_url = 'https://www.pc4u.co.jp/shop/shopbrand.html?search=&prize1=sale'
    new_url = 'https://www.pc4u.co.jp/shop/shopbrand.html?search=new'
    url = "https://www.pc4u.co.jp/shop/shopbrand.html?page=1&search="
    lst = [sale_url, new_url, url]

    crawler_pc4u = CrawlerPc4u(url=outlet_url)
    crawler_pc4u.list_page_crawling()
    for url in lst:
        crawler_pc4u.url = url
        crawler_pc4u.list_page_crawling()
    crawler_pc4u.detail_page_crawling()
    crawler_pc4u.save_excel_file(save_path=settings.SCRAPE_SCHEDULE_SAVE_PATH)
