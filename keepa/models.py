import logging
import threading
import datetime
import os
import configparser

from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import Column, Integer, String, Date
from sqlalchemy.orm import sessionmaker
from sqlalchemy.exc import IntegrityError
from contextlib import contextmanager

config = configparser.ConfigParser()
config.read(os.path.join(os.path.dirname(__file__), 'settings.ini'))
db = config['DB']

logger = logging.getLogger('sqlalchemy.engine')
logger.setLevel(logging.WARNING)
lock = threading.Lock()
postgresql_engine = create_engine(f"postgresql://{db['UserName']}:{db['PassWord']}@{db['Host']}:{db['Port']}/{db['DBname']}")
Base = declarative_base()
NetseaBase = declarative_base()
Session = sessionmaker(bind=postgresql_engine)


class KeepaProducts(Base):
    __tablename__ = 'keepa_products'
    asin = Column(String, primary_key=True)
    sales_drops_90 = Column(Integer)
    created = Column(Date, default=datetime.date.today)
    modified = Column(Date, default=datetime.date.today)

    @classmethod
    def create(cls, asin, drops):
        shop = cls(asin=asin, sales_drops_90=drops)
        try:
            with session_scope() as session:
                session.add(shop)
            return True
        except IntegrityError:
            return False

    @classmethod
    def object_get_db_asin(cls, asin, delay=90):
        with session_scope() as session:
            delay_date = datetime.date.today() - datetime.timedelta(days=delay)
            product = session.query(cls).filter(cls.asin == asin, cls.modified >= delay_date).first()
            if product is None:
                return None
            return product

    @classmethod
    def update_or_insert(cls, asin, drops):
        with session_scope() as session:
            keepa_product = cls(asin=asin, sales_drops_90=drops)
            product = session.query(cls).filter(cls.asin == asin).first()
            if product is None:
                session.add(keepa_product)
                return True
            product.sales_drops_90 = drops
            product.modified = datetime.date.today()
            return True

    @property
    def value(self):
        return {
            'asin': self.asin,
            'sales_drop_90': self.sales_drops_90,
            'created': self.created,
            'modified': self.modified,
        }

@contextmanager
def session_scope():
    session = Session()
    session.expire_on_commit = False
    try:
        lock.acquire()
        yield session
        session.commit()
    except Exception as e:
        logger.error(f'action=session_scope error={e}')
        session.rollback()
        raise
    finally:
        session.expire_on_commit = True
        lock.release()


def init_db():
    Base.metadata.create_all(bind=postgresql_engine)
