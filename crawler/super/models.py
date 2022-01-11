from sqlalchemy import Column
from sqlalchemy import String
from sqlalchemy import BigInteger
from sqlalchemy.sql.expression import and_

from crawler.models import Base
from crawler.models import Product
from crawler.models import session_scope
from crawler.models import Shop
from crawler.models import postgresql_engine


class Super(Product, Base):
    __tablename__ = 'super_products'
    url = Column(String)
    id = Column(BigInteger, primary_key=True, autoincrement=True)

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

class SuperShop(Shop, Base):
    __tablename__ = 'super_shops'


def init_db():
    Base.metadata.create_all(bind=postgresql_engine)
