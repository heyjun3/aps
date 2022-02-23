import time

from bs4 import BeautifulSoup

import log_settings
from crawler.models import Product
from crawler.models import Base
from crawler.models import session_scope
from crawler.models import postgresql_engine
from crawler import utils


logger = log_settings.get_logger(__name__)


class RakutenProduct(Product, Base):
    __tablename__ = 'rakuten_products'


def init_db():
    Base.metadata.create_all(bind=postgresql_engine)