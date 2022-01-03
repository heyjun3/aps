import os
import configparser

from sqlalchemy import create_engine

from crawler.models import Product
from crawler.models import Base


config = configparser.ConfigParser()
config.read(os.path.join(os.path.dirname(__file__), 'settings.ini'))
db = config['DB']

postgresql_engine = create_engine(f"postgresql://{db['UserName']}:{db['PassWord']}@{db['Host']}:{db['Port']}/{db['DBname']}")

class BuffaloProduct(Product, Base):
    __tablename__ = 'buffaloproducts'


def init_db():
    Base.metadata.create_all(bind=postgresql_engine)
