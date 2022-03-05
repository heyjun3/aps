import threading
import datetime

from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import Column, Integer, String, Date
from sqlalchemy import JSON
from sqlalchemy import or_
from sqlalchemy.orm import sessionmaker
from sqlalchemy.orm.attributes import flag_modified
from sqlalchemy.exc import IntegrityError
from contextlib import contextmanager
import pandas as pd
import numpy as np

import settings
import log_settings


logger = log_settings.get_logger(__name__)
lock = threading.Lock()
postgresql_engine = create_engine(settings.DB_URL)
Base = declarative_base()
NetseaBase = declarative_base()
Session = sessionmaker(bind=postgresql_engine)


def convert_render_price_rank_data(context) -> dict|None:
    params = context.get_current_parameters()
    rank_data = params.get('rank_data')
    price_data = params.get('price_data')
    if rank_data is None or price_data is None:
        return None
    rank_dict = {convert_keepa_time_to_datetime_date(int(k)): v for k, v in rank_data.items()}
    price_dict = {convert_keepa_time_to_datetime_date(int(k)): v for k, v in price_data.items()}

    rank_df = pd.DataFrame(data=list(rank_dict.items()), columns=['date', 'rank']).astype({'rank': int})
    price_df = pd.DataFrame(data=list(price_dict.items()), columns=['date', 'price']).astype({'price': int})

    df = pd.merge(rank_df, price_df, on='date', how='outer')
    df = df.replace(-1.0, np.nan)
    df = df.fillna(method='ffill')
    df = df.fillna(method='bfill')
    df = df.replace([np.nan], [None])
    delay = datetime.datetime.now().date() - datetime.timedelta(days=90)
    df = df[df['date'] > delay]
    df = df.sort_values('date', ascending=True)
    products = {'date': list(map(lambda x: x.isoformat(), df['date'].to_list())), 
                'rank': df['rank'].to_list(), 
                'price': df['price'].to_list()}

    return products


class KeepaProducts(Base):
    __tablename__ = 'keepa_products'
    asin = Column(String, primary_key=True)
    sales_drops_90 = Column(Integer)
    created = Column(Date, default=datetime.date.today)
    modified = Column(Date, default=datetime.date.today, onupdate=datetime.date.today)
    price_data = Column(JSON)
    rank_data = Column(JSON)
    render_data = Column(JSON, default=convert_render_price_rank_data, onupdate=convert_render_price_rank_data)

    @classmethod
    def object_get_db_asin(cls, asin, delay=30):
        with session_scope() as session:
            delay_date = datetime.date.today() - datetime.timedelta(days=delay)
            product = session.query(cls).filter(cls.asin == asin, cls.modified >= delay_date).first()
            if product is None:
                return None
            return product

    @classmethod
    def update_or_insert(cls, asin, drops, price_data, rank_data):
        with session_scope() as session:
            product = session.query(cls).filter(cls.asin == asin).first()
            if product is None:
                keepa_product = cls(asin=asin, sales_drops_90=drops, price_data=price_data, rank_data=rank_data)
                try:
                    session.add(keepa_product)
                except IntegrityError as ex:
                    logger.error(ex)
                    logger.error(keepa_product.value)
            else:
                product.sales_drops_90 = drops
                product.price_data = price_data
                product.rank_data = rank_data
            return True

    @classmethod
    def update_price_and_rank_data(cls, asin:str, unix_time: float, price: int, rank: int) -> bool:
        logger.info('action=update_price_and_rank_data status=run')
        time = convert_unix_time_to_keepa_time(unix_time)
        with session_scope() as session:
            product = session.query(cls).filter(cls.asin == asin).first()
            if product is None:
                return False

            product.price_data[time] = price
            product.rank_data[time] = rank
            flag_modified(product, 'price_data')
            flag_modified(product, 'rank_data')

        logger.info('action=update_price_and_rank_data status=done')
        return True

    @classmethod
    def get_asin_list_price_data_is_None(cls, max_count: int = 100):
        with session_scope() as session:
            products = session.query(cls.asin).filter(or_(cls.price_data == None, cls.rank_data == None)).limit(max_count).all()
            if products:
                products = list(map(lambda x: x[0], products))
                return products
            else:
                return None

    @classmethod
    def get_keepa_product(cls, asin: str):
        with session_scope() as session:
            product = session.query(cls).filter(cls.asin == asin).first()
            if product:
                return product
            else:
                return None

    @classmethod
    def get_products_not_modified(cls, count: int=864000):
        today = datetime.date.today()
        with session_scope() as session:
            products = session.query(cls.asin).filter(cls.modified != today, cls.price_data != None, cls.rank_data != None).limit(count).all()
            return [product[0] for product in products]
    
    @classmethod
    def set_render_data(cls):
        with session_scope() as session:
            asin_list = session.query(cls.asin).all()
        for asin in asin_list:
            with session_scope() as session:
                keepa = session.query(cls).filter(cls.asin == asin[0]).first()
                keepa.render_data = keepa.convert_render_price_rank_data()

    @property
    def value(self):
        return {
            'asin': self.asin,
            'sales_drop_90': self.sales_drops_90,
            'created': self.created,
            'modified': self.modified,
            'price_data': self.price_data,
            'rank_data': self.rank_data,
            'render_data': self.render_data,
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


def convert_keepa_time_to_datetime_date(keepa_time: int):
    unix_time = (keepa_time + 21564000) * 60
    date_time = datetime.datetime.fromtimestamp(unix_time)
    return date_time.date()


def convert_unix_time_to_keepa_time(unix_time: float) -> str:
    keepa_time = round(unix_time / 60 - 21564000)
    return str(keepa_time)
