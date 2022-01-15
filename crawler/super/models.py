from logging import getLogger

from sqlalchemy import Column
from sqlalchemy import String
from sqlalchemy import BigInteger
from sqlalchemy import Integer
from sqlalchemy import ForeignKey
from sqlalchemy.exc import IntegrityError
from sqlalchemy.sql.expression import and_

from crawler.models import Base
from crawler.models import session_scope
from crawler.models import Shop
from crawler.models import postgresql_engine


logger = getLogger(__name__)


class Super(Base):
    __tablename__ = 'super_products'
    name = Column(String)
    product_code = Column(String, primary_key=True, nullable=False)
    url = Column(String)
    price = Column(BigInteger)
    shop_code = Column(String)

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

    def save(self):
        try:
            with session_scope() as session:
                session.add(self)
                return True
        except IntegrityError as ex:
            logger.error(ex)



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

    def save(self):
        try:
            with session_scope() as session:
                session.add(self)
                return True
        except IntegrityError as ex:
            logger.error(ex)
            return False

    def save_or_update(self):
        try:
            with session_scope() as session:
                session.add(self)
                return True
        except IntegrityError as ex:
            logger.error(ex)
            with session_scope() as session:
                product = session.query(SuperProductDetails).filter(SuperProductDetails.product_code == self.product_code, SuperProductDetails.set_number == self.set_number).first()
                product.price = self.price
                return True


class SuperShop(Shop, Base):
    __tablename__ = 'super_shops'

    
def init_db():
    Base.metadata.create_all(bind=postgresql_engine)