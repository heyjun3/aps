from contextlib import contextmanager
import datetime
import threading
from copy import deepcopy
import itertools
from typing import List
import asyncio

from sqlalchemy import create_engine
from sqlalchemy import Column
from sqlalchemy import String
from sqlalchemy import Float
from sqlalchemy import DateTime
from sqlalchemy import BigInteger
from sqlalchemy import or_
from sqlalchemy import distinct
from sqlalchemy import Numeric
from sqlalchemy import Computed
from sqlalchemy import func
from sqlalchemy.future import select
from sqlalchemy.orm import sessionmaker
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.ext.asyncio import create_async_engine

from keepa.models import KeepaProducts
import settings
import log_settings


engine = create_engine(settings.DB_URL, pool_pre_ping=True, pool_size=10, connect_args={'connect_timeout': 10})
async_engine = create_async_engine(settings.DB_ASYNC_URL)
Session = sessionmaker(bind=engine)
async_session = sessionmaker(bind=async_engine, class_=AsyncSession)
Base = declarative_base()
lock = threading.Lock()
logger = log_settings.get_logger(__name__)


class NonePriceError(Exception):
    pass


PROFIT = "price - (cost * unit) - ((price * fee_rate) * 1.1) - shipping_fee"
PROFIT_RATE = "(price - (cost * unit) - ((price * fee_rate) * 1.1) - shipping_fee) / price"


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
    profit = Column(BigInteger, Computed(PROFIT))
    profit_rate =Column(Numeric(precision=10, scale=2), Computed(PROFIT_RATE))
    created_at = Column(DateTime, default=datetime.datetime.now)

    def save(self):
        with session_scope() as session:
            try:
                session.add(self)
            except IntegrityError as ex:
                logger.debug(ex)
                return False
            return True

    @classmethod
    def get_filenames(cls):
        with session_scope() as session:
            filenames = session.query(distinct(cls.filename)).all()
            filenames = itertools.chain.from_iterable(filenames)

        return sorted(filenames)

    @classmethod
    def get_done_filenames(cls):
        with session_scope() as session:
            filenames = session.query(distinct(cls.filename)).all()
        filenames = itertools.chain.from_iterable(filenames) 

        with session_scope() as session:
            not_done_filenames = session.query(distinct(cls.filename)).join(KeepaProducts, cls.asin == KeepaProducts.asin, isouter=True)\
                                 .filter(cls.profit >= 200, cls.profit_rate >= 0.1, KeepaProducts.asin == None).all()
        not_done_filenames = itertools.chain.from_iterable(not_done_filenames)

        return sorted(list(set(filenames) - set(not_done_filenames)))

    @classmethod
    def get_render_data(cls, filename: str, page: int=1, count: int=500):
        start = (page - 1) * count
        end = start + count
        with session_scope() as session:
            rows = session.query(cls, KeepaProducts.render_data).join(KeepaProducts, cls.asin == KeepaProducts.asin)\
            .filter(cls.profit >= 200, cls.profit_rate >= 0.1, cls.filename == filename, 
            cls.unit <= 10, KeepaProducts.sales_drops_90 > 3, KeepaProducts.render_data != None)\
            .order_by(KeepaProducts.sales_drops_90.desc()).slice(start, end).all()
            return rows

    @classmethod
    def get_chart_data(cls, filename: str, page: int, count: int) -> List: 
        start = (page - 1) * count
        end = start + count
        with session_scope() as session:
            rows = session.query(cls, KeepaProducts.render_data).join(KeepaProducts, cls.asin == KeepaProducts.asin)\
                    .filter(cls.profit >= 200, 
                            cls.profit_rate >= 0.1, 
                            cls.filename == filename, 
                            cls.unit <= 10, 
                            KeepaProducts.sales_drops_90 > 3, 
                            KeepaProducts.render_data != None)\
                    .order_by(MWS.profit.desc()).slice(start, end).all()
        return rows

    @classmethod
    def get_row_count(cls, filename: str) -> int:
        with session_scope() as session:
            rows = session.query(func.count(cls.asin)).join(KeepaProducts, cls.asin == KeepaProducts.asin)\
            .filter(cls.profit >= 200, cls.profit_rate >= 0.1, cls.filename == filename, cls.unit <= 10, KeepaProducts.sales_drops_90 > 3, KeepaProducts.render_data != None).first()
            return rows[0]

    @classmethod
    def get_asin_list_None_products(cls, profit: int=200, profit_rate: float=0.1):
        with session_scope() as session:
            try:
                asin_list = session.query(cls.asin).filter(cls.profit >= profit, cls.profit_rate >= profit_rate)\
                            .join(KeepaProducts, cls.asin == KeepaProducts.asin, isouter=True).filter(KeepaProducts.asin == None).all()
                asin_list = list(set(map(lambda x: x[0], asin_list)))
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
    def get_price_is_None_asins(cls):
        with session_scope() as session:
            products = session.query(cls.asin).filter(cls.price == None).all()
            products = list(map(lambda x: x[0], products))
            return products

    @classmethod
    def get_fee_is_None_products(cls):
        with session_scope() as session:
            products = session.query(cls.asin).filter(or_(cls.fee_rate == None, cls.shipping_fee == None)).all()
            return products

    @classmethod
    def get_fee_is_None_asins(cls, limit_count: int=1000):
        with session_scope() as session:
            products = session.query(cls.asin).filter(or_(cls.fee_rate == None, cls.shipping_fee == None))\
            .order_by(cls.created_at).limit(limit_count).all()
            products = list(map(lambda x: x[0], products))
            return products

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

    @classmethod
    def delete_rows(cls, filename: str):
        with session_scope() as session:
            session.query(cls).filter(cls.filename == filename).delete()
        return True

    @classmethod
    def delete_rows_lower_price(cls, profit: int=200, profit_rate: float=0.1, unit_count: int=10, drops: int=3) -> bool:
        with session_scope() as session:
            products = session.query(cls.filename, cls.asin).join(KeepaProducts, cls.asin == KeepaProducts.asin).where(or_(
                cls.profit < profit,
                cls.profit_rate < profit_rate,
                cls.unit > unit_count,
                KeepaProducts.sales_drops_90 <= drops,
            )).all()
            for filename, asin in products:
                session.query(cls).where(cls.filename == filename, cls.asin == asin).delete()

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
    
    def __repr__(self):
        return f'{self.__class__}'


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


def init_db(engine):
    Base.metadata.create_all(bind=engine)


async def async_main():
    async with async_session() as session:
        stmt = select(MWS).where(MWS.asin=='B000FQRD6S').limit(1)
        result = await session.execute(stmt)
        product = result.scalars().first()
        print(product.value)
        product.unit = 3
        await session.commit()
