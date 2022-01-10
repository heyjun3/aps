from logging import getLogger
import time

from bs4 import BeautifulSoup

from crawler.models import Product
from crawler.models import Base
from crawler.models import session_scope
from crawler.models import postgresql_engine
from crawler import utils


logger = getLogger(__name__)


class RakutenProduct(Product, Base):
    __tablename__ = 'rakuten_products'

    def get_jan_code(self):
        logger.info('action=get_jan_code status=run')
        """self has jan code"""
        if self.jan:
            return True

        jan_code = self.fetch()

        if jan_code is None:
            response = utils.request(self.url)
            time.sleep(2)
            self.jan = self.scraping_jan_code(response.text)
            self.save()
        else:
            self.jan = jan_code

    @staticmethod
    def scraping_jan_code(response: str):
        soup = BeautifulSoup(response, 'lxml')
        try:
            jan = soup.select_one('#ratRanCode').get('value')
        except AttributeError as e:
            logger.error(f'{e}')
            return None
        return jan

    def fetch(self):
        with session_scope() as session:
            product = session.query(RakutenProduct).filter(RakutenProduct.url == self.url).first()
            if not product:
                return None
            elif not product.jan:
                return None
            elif not product.price == self.price:
                product.price = self.price
                return product.jan
            else:
                return product.jan

def init_db():
    Base.metadata.create_all(bind=postgresql_engine)