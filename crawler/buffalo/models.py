import os
import configparser

from sqlalchemy import create_engine

from crawler.models import Product
from crawler.models import Base
from crawler.models import postgresql_engine
import settings


class BuffaloProduct(Product, Base):
    __tablename__ = 'buffalo_products'


def init_db():
    Base.metadata.create_all(bind=postgresql_engine)
