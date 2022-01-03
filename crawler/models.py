import logging
import threading
import datetime
import time
import os
import configparser

from bs4 import BeautifulSoup
from sqlalchemy import create_engine, Index
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import Column, Integer, String, Date, Float
from sqlalchemy.orm import sessionmaker
from sqlalchemy.exc import IntegrityError
from sqlalchemy.sql.expression import and_
from contextlib import contextmanager

import utils


logger = logging.getLogger('sqlalchemy.engine')
logger.setLevel(logging.WARNING)
config = configparser.ConfigParser()
config.read(os.path.join(os.path.dirname(__file__), 'settings.ini'))
db = config['DB']

lock = threading.Lock()
postgresql_engine = create_engine(f"postgresql://{db['UserName']}:{db['PassWord']}@{db['Host']}:{db['Port']}/{db['DBname']}")
Base = declarative_base()
Session = sessionmaker(bind=postgresql_engine)


class NetseaShopUrl(Base):
    __tablename__ = 'NetseaShopUrl'
    id = Column(Integer, primary_key=True, autoincrement=True)
    url = Column(String)
    shop_id = Column(Integer)
    quantity = Column(Integer)

    def save(self):
        try:
            with session_scope() as session:
                session.add(self)
            return True
        except IntegrityError:
            return False


class KeepaProducts(Base):
    __tablename__ = 'KeepaProducts'
    asin = Column(String, primary_key=True)
    sales_drops_90 = Column(Integer)
    created = Column(Date, default=datetime.date.today)
    modified = Column(Date, default=datetime.date.today)

    @classmethod
    def create(cls, asin, drops):
        shop = cls(asin=asin, sales_drops_90=drops)
        try:
            with session_scope() as session:
                session.add(shop)
            return True
        except IntegrityError:
            return False

    @classmethod
    def object_get_db_asin(cls, asin, delay=90):
        with session_scope() as session:
            delay_date = datetime.date.today() - datetime.timedelta(days=delay)
            product = session.query(cls).filter(cls.asin == asin, cls.modified >= delay_date).first()
            if product is None:
                return None
            return product

    @classmethod
    def update_or_insert(cls, asin, drops):
        with session_scope() as session:
            keepa_product = cls(asin=asin, sales_drops_90=drops)
            product = session.query(cls).filter(cls.asin == asin).first()
            if product is None:
                session.add(keepa_product)
                return True
            product.sales_drops_90 = drops
            product.modified = datetime.date.today()
            return True

    @property
    def value(self):
        return {
            'asin': self.asin,
            'sales_drop_90': self.sales_drops_90,
            'created': self.created,
            'modified': self.modified,
        }


class Shop:
    name = Column(String)
    shop_id = Column(Integer, primary_key=True, nullable=False)
    url = Column(String)
    product_quantity = Column(Integer)
    category_id = Column(Integer)

    @classmethod
    def create(cls, name, shop_id, url, quantity, category):
        shop = cls(name=name, shop_id=shop_id, url=url, product_quantity=quantity, category_id=category)
        try:
            with session_scope() as session:
                session.add(shop)
            return True
        except IntegrityError:
            return False

    @classmethod
    def shop_id_get_shop_info(cls, shop_id):
        with session_scope() as session:
            shop_info = session.query(cls).filter(cls.shop_id == shop_id).first()
            if shop_info is None:
                return None
            return shop_info

    @classmethod
    def get_all_info(cls):
        with session_scope() as session:
            shops = session.query(cls).all()
            return shops

    @property
    def value(self):
        return {
            'name': self.name,
            'shop_id': self.shop_id,
            'url': self.url,
            'product_quantity': self.product_quantity,
            'category_id': self.category_id,
        }

    def save(self):
        try:
            with session_scope() as session:
                session.add(self)
            return True
        except IntegrityError:
            return False


class Product:
    name = Column(String)
    jan = Column(String)
    price = Column(Integer)
    shop_code = Column(String)
    url = Column(String, primary_key=True, nullable=False)
    product_code = Column(String)

    @classmethod
    def create(cls, name=None, jan=None, price=None, shop_code=None, url=None, product_code=None):
        product = cls()
        product.name = name
        product.jan = jan
        product.price = price
        product.shop_code = shop_code
        product.url = url
        product.product_code = product_code
        return product

    @classmethod
    def get(cls, url):
        with session_scope() as session:
            product = session.query(cls).filter(cls.url == url).first()
            if product is None:
                return None
            return product

    @classmethod
    def get_jan(cls, jan):
        with session_scope() as session:
            product = session.query(cls).filter(cls.jan == jan).first()
            if product is None:
                return None
            return product

    @classmethod
    def get_all_info(cls):
        with session_scope() as session:
            products = session.query(cls).all()
            return products

    @classmethod
    def get_object_filter_productcode_shopcode(cls, product_code, shop_code):
        print('action=get_object_filter_productcode_shopcode status=run')
        with session_scope() as session:
            product = session.query(cls).filter(and_(cls.product_code == product_code, cls.shop_code == shop_code)).first()
            if product is None:
                return None
            return product

    @classmethod
    def price_update(cls, select_col, select_value, new_value):
        with session_scope() as session:
            if select_col == 'product_code':
                product = session.query(cls).filter(cls.product_code == select_value).first()
                product.price = new_value
            return True

    @classmethod
    def all_update(cls):
        with session_scope() as session:
            products = session.query(cls).all()
            for product in products:
                product.price = int(product.price)

            return True

    def save(self):
        try:
            with session_scope() as session:
                session.add(self)
            return True
        except IntegrityError:
            return False

    @classmethod
    def get_product_code(cls, product_code):
        with session_scope() as session:
            products = session.query(cls).filter(cls.product_code == product_code).first()
            if products is None:
                return None
            return products

    @property
    def value(self):
        return {
            'name': self.name,
            'jan': self.jan,
            'price': self.price,
            'shop_code': self.shop_code,
            'url': self.url,
            'product_code': self.product_code,
        }


class Netsea(Product, Base):
    __tablename__ = 'Netseaproducts'


class Super(Product, Base):
    __tablename__ = 'Superproducts'
    url = Column(String)
    id = Column(Integer, primary_key=True, autoincrement=True)

    @classmethod
    def get_product(cls, product_code, price):
        with session_scope() as session:
            products = session.query(cls).filter(and_(cls.product_code == product_code, cls.price == price)).all()
            if not products:
                return None
            products = session.query(cls).filter(cls.product_code == product_code).all()
            return products

    @classmethod
    def get_product_jan_and_update_price(cls, product_code, jan, price):
        with session_scope() as session:
            product = session.query(cls).filter(cls.product_code == product_code, cls.jan == jan).first()
            if not product:
                return None
            product.price = price
            return True

    @classmethod
    def get_url(cls, url):
        with session_scope() as session:
            products = session.query(cls).filter(cls.url == url, cls.jan.isnot(None)).all()
            if not products:
                return None
            return products


class Pc4u(Product, Base):
    __tablename__ = 'pc4uproducts'


class RakutenProduct(Product, Base):
    __tablename__ = 'rakutenproducts'

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


class Test(Product, Base):
    __tablename__ = 'test'


class NetseaShop(Shop, Base):
    __tablename__ = 'NetseaShops'
    discount_rate = Column(Float)


class SuperShop(Shop, Base):
    __tablename__ = 'SuperShops'


class TestShop(Shop, Base):
    __tablename__ = 'testshop'


@contextmanager
def session_scope():
    session = Session()
    session.expire_on_commit = False
    try:
        lock.acquire()
        yield session
        session.commit()
    except Exception as e:
        logger.error(f'action=session_scope error={e}')
        session.rollback()
        raise
    finally:
        session.expire_on_commit = True
        lock.release()


def init_db():
    Base.metadata.create_all(bind=postgresql_engine)


# if __name__ == '__main__':
    # transfer_db()
    # init_db()
    # session = Session()
    # shops = session.query(NetseaShop).all()
    # print(shops)

    # shop_info = TestShop.get_all_info()
    # for shop in shop_info:
    #     NetseaShop.create(name=shop.name, shop_id=shop.shop_id, url=shop.url,
    #                       quantity=shop.product_quantity, category=shop.category_id)

    # with open('tmp', 'r') as f:
    #     for i in f.readlines():
    #         print(i.strip(), type(i))
    #         shop_id = int(i.strip())
    #         with session_scope() as session:
    #             shop = session.query(NetseaShop).filter(NetseaShop.shop_id == shop_id).first()
    #             if shop:
    #                 shop.discount_rate = 0.95

    # product = Netsea.get_object_filter_productcode_shopcode("49795462", "357136")
    # print(product.value)

    # netseaindex = Index('ix_pc4uproducts_product_code', Pc4u.product_code)
    # netseaindex.create(bind=engine)

