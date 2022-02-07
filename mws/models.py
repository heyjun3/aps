from contextlib import contextmanager
import datetime
from logging import getLogger
import threading

from sqlalchemy import create_engine
from sqlalchemy import Column
from sqlalchemy import String
from sqlalchemy import Integer
from sqlalchemy import Float
from sqlalchemy import or_
from sqlalchemy import distinct
from sqlalchemy.orm import sessionmaker
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.declarative import declarative_base

from keepa.models import KeepaProducts

import settings

engine = create_engine(settings.DB_URL, pool_pre_ping=True)
Session = sessionmaker(bind=engine)
Base = declarative_base()
lock = threading.Lock()
logger = getLogger(__name__)


class MWS(Base):
    __tablename__ = 'mws_products'
    asin = Column(String, primary_key=True, nullable=False)
    filename = Column(String, primary_key=True, nullable=False)
    title = Column(String)
    jan = Column(String)
    unit = Column(Integer)
    price = Column(Integer)
    cost = Column(Integer)
    fee_rate = Column(Float)
    shipping_fee = Column(Integer)

    def save(self):
        with session_scope() as session:
            try:
                session.add(self)
            except IntegrityError as ex:
                logger.debug(ex)
                return False
            return True

    @classmethod
    def get(cls, asin):
        with session_scope() as session:
            profit = (cls.price - (cls.cost * cls.unit) - ((cls.price * cls.fee_rate) * 1.1) - cls.shipping_fee)
            profit_rate = profit / cls.price
            mws = session.query(cls, profit, profit_rate).filter(cls.asin == asin).first()
            if mws:
                return mws
            else:
                return None

    @classmethod
    def get_completion_filename_list(cls):
        with session_scope() as session:
            profit = (cls.price - (cls.cost * cls.unit) - ((cls.price * cls.fee_rate) * 1.1) - cls.shipping_fee)
            profit_rate = profit / cls.price

            mws_sub_query = session.query(distinct(cls.filename)).filter(or_(cls.price == None, cls.fee_rate == None))
            keepa_sub_query = session.query(distinct(cls.filename)).filter(profit > 200, profit_rate > 0.1)\
                              .join(KeepaProducts, cls.asin == KeepaProducts.asin, isouter=True)\
                              .filter(or_(KeepaProducts.asin == None, KeepaProducts.rank_data == None, KeepaProducts.price_data == None))
            
            filename_list = session.query(distinct(cls.filename)).filter(cls.filename.notin_(mws_sub_query.union(keepa_sub_query))).all()
            filename_list = sorted(list(map(lambda x: x[0], filename_list)), key=lambda x: x)

            return filename_list

    @classmethod
    def get_render_data(cls, filename: str):
        with session_scope() as session:
            profit = (cls.price - (cls.cost * cls.unit) - ((cls.price * cls.fee_rate) * 1.1) - cls.shipping_fee)
            profit_rate = profit / cls.price

            rows = session.query(cls, KeepaProducts).filter(profit >= 200, profit_rate >= 0.1, cls.filename == filename)\
                   .join(KeepaProducts, cls.asin == KeepaProducts.asin).filter(KeepaProducts.sales_drops_90 > 3).all()
            return rows

    @classmethod
    def get_asin_to_request_keepa(cls, term=30):
        with session_scope() as session:
            profit = (cls.price - (cls.cost * cls.unit) - ((cls.price * cls.fee_rate) * 1.1) - cls.shipping_fee)
            profit_rate = profit / cls.price
            past_date = datetime.date.today() - datetime.timedelta(days=term)
            try:
                asin_list = session.query(cls.asin).filter(profit > 200, profit_rate > 0.1)\
                            .join(KeepaProducts, cls.asin == KeepaProducts.asin, isouter=True)\
                            .filter(or_(KeepaProducts.asin is None, KeepaProducts.modified < past_date)).all()
                asin_list = list(map(lambda x: x[0], asin_list))
            except Exception as ex:
                logger.error(f'action=get_asin_to_request_keepa error={ex}')
                return None
            return asin_list

    @classmethod
    def get_price_is_None_products(cls):
        with session_scope() as session:
            products = session.query(cls).filter(cls.price == None).all()
            return products

    @classmethod
    def update_price(cls, asin: str, filename: str, price: int):
        with session_scope() as session:
            mws = session.query(cls).filter(cls.asin == asin, cls.filename == filename).first()
            try:
                mws.price = price
            except Exception as ex:
                logger.error(f'action=update_price error={ex}')
                return False
            return True

    @classmethod
    def update_fee(cls, asin: str, filename: str, fee_rate: float, shipping_fee: int):
        with session_scope() as session:
            mws = session.query(cls).filter(cls.asin == asin, cls.filename == filename).first()
            try:
                mws.fee_rate = fee_rate
                mws.shipping_fee = shipping_fee
            except Exception as ex:
                logger.error(f'action=update_fee_and_profit error={ex}')
                return False
            return True

    @classmethod
    def delete_objects(cls, filename: str):
        with session_scope() as session:
            session.query(cls).filter(cls.filename == filename).delete()
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

    @property
    def profit(self):
        return int(self.price - (self.cost * self.unit) - ((self.price * self.fee_rate) * 1.1) - self.shipping_fee)

    @property
    def profit_rate(self):
        return round(self.profit / self.price, 2)


@contextmanager
def session_scope():
    session = Session()
    session.expire_on_commit = False
    try:
        lock.acquire()
        yield session
        session.commit()
    except Exception as ex:
        logger.error(f'action=session_scope error={ex}')
        session.rollback()
    finally:
        session.expire_on_commit = True
        lock.release()


def init_db():
    Base.metadata.create_all(bind=engine)
