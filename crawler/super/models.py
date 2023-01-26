from sqlalchemy import Column
from sqlalchemy import String
from sqlalchemy import BigInteger
from sqlalchemy import Integer
from sqlalchemy import ForeignKey
from sqlalchemy.exc import IntegrityError
from sqlalchemy.sql.expression import and_

import log_settings
from crawler.models import Base
from crawler.models import session_scope
from crawler.models import Shop
from crawler.models import postgresql_engine


logger = log_settings.get_logger(__name__)


class SuperProduct(Base):
    __tablename__ = 'super_products'
    product_code = Column(String, primary_key=True, nullable=False)
    name = Column(String)
    price = Column(BigInteger)
    shop_code = Column(String)
    url = Column(String)

    def __init__(self, product_code, name=None, price=None, shop_code=None, url=None):
        self.product_code = product_code
        self.name = name
        self.price = price
        self.shop_code = shop_code
        self.url = url

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

    def save(self):
        try:
            with session_scope() as session:
                session.add(self)
                return True
        except IntegrityError as ex:
            logger.info(ex)
            return False
        except Exception as ex:
            logger.error(ex)
            return False



class SuperProductDetails(Base):
    __tablename__ = 'super_product_details'
    product_code = Column(String, ForeignKey("super_products.product_code"), nullable=False, primary_key=True)
    set_number = Column(Integer, primary_key=True)
    shop_code = Column(String)
    price = Column(BigInteger)
    jan = Column(String)

    @property
    def value(self):
        return {
            'product_code': self.product_code,
            'set_number': self.set_number,
            'shop_code': self.shop_code,
            'price': self.price,
            'jan': self.jan,
        }

    @classmethod
    def get(cls, product_code, price):
        with session_scope() as session:
            products = session.query(cls).filter(cls.product_code == product_code).all()
            for product in products:
                if product.price == price:
                    return products
            return None

    @classmethod
    def get_objects_to_product_code(cls, product_code: str):
        with session_scope() as session:
            products = session.query(cls).filter(cls.product_code == product_code).all()
            if products:
                return products
            else:
                return None
                
    def save(self):
        try:
            with session_scope() as session:
                session.add(self)
                return True
        except IntegrityError as ex:
            logger.info(ex)
            return False
        except Exception as ex:
            logger.error(ex)
            return False

    def save_or_update(self):
        try:
            with session_scope() as session:
                session.add(self)
                return True
        except IntegrityError as ex:
            logger.info(ex)
            with session_scope() as session:
                product = session.query(SuperProductDetails).filter(SuperProductDetails.product_code == self.product_code, SuperProductDetails.set_number == self.set_number).first()
                product.price = self.price
                return True
        except Exception as ex:
            logger.error(ex)
            return False


class SuperShop(Shop, Base):
    __tablename__ = 'super_shops'

    
def init_db():
    Base.metadata.create_all(bind=postgresql_engine)
