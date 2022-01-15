from logging import getLogger

from sqlalchemy import Column
from sqlalchemy import Integer
from sqlalchemy import String
from sqlalchemy import Float
from sqlalchemy.exc import IntegrityError

from crawler.models import session_scope
from crawler.models import Base
from crawler.models import postgresql_engine
from crawler.models import Product
from crawler.models import Shop


logger = getLogger(__name__)


class Netsea(Product, Base):
    __tablename__ = 'netsea_products'


class NetseaShop(Shop, Base):
    __tablename__ = 'netsea_shops'
    discount_rate = Column(Float)


class NetseaShopUrl(Base):
    __tablename__ = 'netsea_shop_url'
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


def init_db():
    Base.metadata.create_all(bind=postgresql_engine)
