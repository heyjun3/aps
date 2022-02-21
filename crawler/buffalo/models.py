from sqlalchemy import Column
from sqlalchemy import String

from crawler.models import Product
from crawler.models import Base
from crawler.models import postgresql_engine


class BuffaloProduct(Product, Base):
    __tablename__ = 'buffalo_products'

    shop_code = Column(String, primary_key=True, nullable=False, default='buffalo')


def init_db():
    Base.metadata.create_all(bind=postgresql_engine)
