from contextlib import contextmanager
import datetime
from logging import getLogger
import threading
from copy import deepcopy

from sqlalchemy import create_engine
from sqlalchemy import Column
from sqlalchemy import String
from sqlalchemy import Float
from sqlalchemy import BigInteger
from sqlalchemy import or_
from sqlalchemy import distinct
from sqlalchemy.orm import sessionmaker
from sqlalchemy.orm import validates
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.declarative import declarative_base

from keepa.models import KeepaProducts
import settings
import log_settings


engine = create_engine(settings.DB_URL, pool_pre_ping=True, pool_size=10, connect_args={'connect_timeout': 10})
Session = sessionmaker(bind=engine)
Base = declarative_base()
lock = threading.Lock()
logger = log_settings.get_logger(__name__)



def calc_profit(context) -> int|None:
    params = context.get_current_parameters()
    print(params)
    try:
        return int(params['price'] - (params['cost'] * params['unit']) - ((params['price'] * params['fee_rate']) * 1.1) - params['shipping_fee'])
    except TypeError as ex:
        logger.error(ex)
        return None

def calc_profit_rate(context) -> float|None:
    params = context.get_current_parameters()
    try:
        return round(params['profit'] / params['price'], 2)
    except TypeError as ex:
        logger.error(ex)
        return None


class NonePriceError(Exception):
    pass


class MWS(Base):
    __tablename__ = 'mws_products'
    asin = Column(String, primary_key=True, nullable=False)
    filename = Column(String, primary_key=True, nullable=False)
    title = Column(String)
    jan = Column(String)
    unit = Column(BigInteger)
    price = Column(BigInteger)
    cost = Column(BigInteger)
    fee_rate = Column(Float)
    shipping_fee = Column(BigInteger)
    profit = Column(BigInteger, default=calc_profit, onupdate=calc_profit)
    profit_rate = Column(Float, default=calc_profit_rate, onupdate=calc_profit_rate)

    def save(self):
        with session_scope() as session:
            try:
                session.add(self)
            except IntegrityError as ex:
                logger.debug(ex)
                return False
            return True

    @classmethod
    def get_completion_filename_list(cls):
        with session_scope() as session:

            mws_sub_query = session.query(distinct(cls.filename)).filter(or_(cls.price == None, cls.fee_rate == None))
            keepa_sub_query = session.query(distinct(cls.filename)).filter(cls.profit > 200, cls.profit_rate > 0.1)\
                              .join(KeepaProducts, cls.asin == KeepaProducts.asin, isouter=True)\
                              .filter(or_(KeepaProducts.asin == None, KeepaProducts.rank_data == None, KeepaProducts.price_data == None))
            
            filename_list = session.query(distinct(cls.filename)).filter(cls.filename.notin_(mws_sub_query.union(keepa_sub_query))).all()
            filename_list = sorted(list(map(lambda x: x[0], filename_list)), key=lambda x: x)

            return filename_list

    @classmethod
    def get_render_data(cls, filename: str):
        with session_scope() as session:
            rows = session.query(cls, KeepaProducts.render_data).filter(cls.profit >= 200, cls.profit_rate >= 0.1, cls.filename == filename)\
                   .join(KeepaProducts, cls.asin == KeepaProducts.asin).filter(KeepaProducts.sales_drops_90 > 3).all()
            return rows

    @classmethod
    def get_asin_list_None_products(cls, term=30, count=100):
        with session_scope() as session:
            past_date = datetime.date.today() - datetime.timedelta(days=term)
            try:
                asin_list = session.query(cls.asin).filter(cls.profit > 200, cls.profit_rate > 0.1)\
                            .join(KeepaProducts, cls.asin == KeepaProducts.asin, isouter=True)\
                            .filter(or_(KeepaProducts.asin == None, KeepaProducts.modified < past_date, KeepaProducts.rank_data == None, KeepaProducts.price_data == None))\
                            .limit(count).all()
                asin_list = list(map(lambda x: x[0], asin_list))
            except Exception as ex:
                logger.error(f'action=get_asin_to_request_keepa error={ex}')
                return None
            return asin_list

    @classmethod
    def get_price_is_None_products(cls):
        with session_scope() as session:
            products = session.query(cls.asin).filter(cls.price == None).all()
            return products

    @classmethod
    def get_fee_is_None_products(cls):
        with session_scope() as session:
            products = session.query(cls.asin).filter(or_(cls.fee_rate == None, cls.shipping_fee == None)).all()
            return products

    @validates('cost')
    @classmethod
    def update_price(cls, asin: str, price: int):
        with session_scope() as session:
            mws_list = session.query(cls).filter(cls.asin == asin).all()
            for mws in mws_list:
                mws.price = price
                mws.cost = deepcopy(mws.cost)
            return True

    @classmethod
    def update_fee(cls, asin: str, fee_rate: float, shipping_fee: int):
        with session_scope() as session:
            mws_list = session.query(cls).filter(cls.asin == asin).all()
            for mws in mws_list:
                mws.fee_rate = fee_rate
                mws.shipping_fee = shipping_fee
            return True

    @classmethod
    def delete_objects(cls, filename_list: list):
        with session_scope() as session:
            session.query(cls).filter(cls.filename.in_(filename_list)).delete()
            return True

    @property
    def value(self):
        return {
            'asin': self.asin,
            'filename': self.filename,
            'jan': self.jan,
            'unit': self.unit,
            'cost': self.cost,
            'title': self.title,
            'price': self.price,
            'fee_rate': self.fee_rate,
            'shipping_fee': self.shipping_fee,
            'profit': self.profit,
            'profit_rate': self.profit_rate,
        }


@contextmanager
def session_scope():
    session = Session()
    session.expire_on_commit = False
    try:
        lock.acquire()
        yield session
        session.commit()
    except IntegrityError as ex:
        logger.debug(f'action=session_scope error={ex}')
        session.rollback()
    except Exception as ex:
        logger.error(f'action=session_scope error={ex}')
        session.rollback()
    finally:
        session.expire_on_commit = True
        lock.release()


def init_db():
    Base.metadata.create_all(bind=engine)
