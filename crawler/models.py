import threading

from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import Column, Integer, String, Float, BigInteger
from sqlalchemy.orm import sessionmaker
from sqlalchemy.exc import IntegrityError
from sqlalchemy.sql.expression import and_
from contextlib import contextmanager

import settings
import log_settings


logger = log_settings.get_logger(__name__)

lock = threading.Lock()
postgresql_engine = create_engine(settings.DB_URL)
Base = declarative_base()
Session = sessionmaker(bind=postgresql_engine)


class Shop:
    name = Column(String)
    shop_id = Column(String, primary_key=True, nullable=False)

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
        }

    def save(self):
        try:
            with session_scope() as session:
                session.add(self)
            return True
        except IntegrityError:
            return False

    @classmethod
    def delete(cls):
        with session_scope() as session:
            session.query(cls).delete()


class Product:
    name = Column(String)
    jan = Column(String)
    price = Column(BigInteger)
    shop_code = Column(String, primary_key=True, nullable=False)
    product_code = Column(String, primary_key=True, nullable=False)

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
    def get_object_filter_productcode_and_shopcode(cls, product_code, shop_code):
        logger.info('action=get_object_filter_productcode_and_shopcode status=run')
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
            'product_code': self.product_code,
        }


@contextmanager
def session_scope():
    session = Session()
    session.expire_on_commit = False
    try:
        lock.acquire()
        yield session
        session.commit()
    except IntegrityError as e:
        logger.debug(f'action=session_scope error={e}')
        session.rollback()
    except Exception as e:
        logger.error(f'action=session_scope error={e}')
        session.rollback()
    finally:
        session.expire_on_commit = True
        lock.release()


def init_db():
    Base.metadata.create_all(bind=postgresql_engine)
