from sqlalchemy import Column
from sqlalchemy import Integer
from sqlalchemy import String
from sqlalchemy import Float
from sqlalchemy.exc import IntegrityError

import log_settings
from crawler.models import session_scope
from crawler.models import Base
from crawler.models import postgresql_engine
from crawler.models import Product
from crawler.models import Shop


logger = log_settings.get_logger(__name__)


class NetseaProduct(Product, Base):
    __tablename__ = 'netsea_products'


class NetseaShop(Shop, Base):
    __tablename__ = 'netsea_shops'


def init_db():
    Base.metadata.create_all(bind=postgresql_engine)
