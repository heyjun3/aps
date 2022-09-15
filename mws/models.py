from __future__ import annotations
import datetime
import threading
from typing import List

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
from sqlalchemy import update
from sqlalchemy import delete
from sqlalchemy.future import select
from sqlalchemy.orm import sessionmaker
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.ext.asyncio import create_async_engine

from models_base import ModelsBase
from keepa.models import KeepaProducts
import settings
import log_settings


engine = create_engine(settings.DB_URL, pool_pre_ping=True, pool_size=10, connect_args={'connect_timeout': 10})
async_engine = create_async_engine(settings.DB_ASYNC_URL)
Session = sessionmaker(bind=engine)
session_scope = sessionmaker(bind=async_engine, class_=AsyncSession)
Base = declarative_base()
lock = threading.Lock()
logger = log_settings.get_logger(__name__)


class NonePriceError(Exception):
    pass


PROFIT = "price - (cost * unit) - ((price * fee_rate) * 1.1) - shipping_fee"
PROFIT_RATE = "(price - (cost * unit) - ((price * fee_rate) * 1.1) - shipping_fee) / price"


class MWS(Base, ModelsBase):
    __tablename__ = 'mws_products'
    asin = Column(String, primary_key=True, nullable=False)
    filename = Column(String, primary_key=True, nullable=False)
    title = Column(String)
    jan = Column(String)
    unit = Column(BigInteger)
    url = Column(String)
    price = Column(BigInteger)
    cost = Column(BigInteger)
    fee_rate = Column(Float)
    shipping_fee = Column(BigInteger)
    profit = Column(BigInteger, Computed(PROFIT))
    profit_rate =Column(Numeric(precision=10, scale=2), Computed(PROFIT_RATE))
    created_at = Column(DateTime, default=datetime.datetime.now)

    async def save(self):
        async with self.session_scope() as session:
            session.add(self)
        return True

    @classmethod
    async def get(cls, asin: str) -> MWS:
        async with cls.session_scope() as session:
            stmt = select(cls).where(cls.asin == asin)
            result = await session.execute(stmt)
            return result.scalar()

    @classmethod
    async def get_filenames(cls) -> List[str]:
        async with cls.session_scope() as session:
            stmt = select(distinct(cls.filename))
            result = await session.execute(stmt)
            filenames = result.scalars()
            return sorted(filenames)

    @classmethod
    async def get_chart_data(cls, filename: str, page: int, count: int) -> List[MWS, dict]:
        start = (page - 1) * count
        end = start + count
        async with cls.session_scope() as session:
            stmt = select(cls, KeepaProducts.render_data).join(KeepaProducts, cls.asin == KeepaProducts.asin).where(
                cls.profit >= 200,
                cls.profit_rate >= 0.1,
                cls.filename == filename,
                cls.unit <= 10,
                KeepaProducts.sales_drops_90 > 3,
                KeepaProducts.render_data != None).order_by(cls.profit.desc()).slice(start, end)
            result = await session.execute(stmt)
            return result.all()

    @classmethod
    async def get_row_count(cls, filename: str) -> int:
        async with cls.session_scope() as session:
            stmt = select(func.count(cls.asin)).join(KeepaProducts, cls.asin == KeepaProducts.asin).where(
                cls.profit >= 200,
                cls.profit_rate >= 0.1,
                cls.filename == filename,
                cls.unit <= 10,
                KeepaProducts.sales_drops_90 > 3,
                KeepaProducts.render_data != None,
            )
            result = await session.execute(stmt)
            return result.scalar()
    @classmethod
    async def get_asin_list_None_products(cls, profit: int=200, profit_rate: float=0.1) -> List[str]:
        async with cls.session_scope() as session:
            stmt = select(cls.asin).join(KeepaProducts, cls.asin == KeepaProducts.asin, isouter=True).where(
                cls.profit >= profit,
                cls.profit_rate >= profit_rate,
                KeepaProducts.asin == None)
            result = await session.execute(stmt)
            return result.scalars().all()

    @classmethod
    async def get_price_is_None_asins(cls) -> List[str]:
        async with cls.session_scope() as session:
            stmt = select(cls.asin).where(cls.price == None)
            result = await session.execute(stmt)
            return result.scalars().all()

    @classmethod
    async def get_fee_is_None_asins(cls, limit_count=10000) -> List[str]:
        async with cls.session_scope() as session:
            stmt = select(cls.asin).where(or_(cls.fee_rate == None, cls.shipping_fee == None))\
                    .order_by(cls.created_at).limit(limit_count)
            result = await session.execute(stmt)
            return result.scalars().all()

    @classmethod
    async def update_price(cls, asin:str, price: int):
        async with cls.session_scope() as session:
            stmt = update(cls).where(cls.asin == asin).values(price=price)
            result = await session.execute(stmt)
            result.close()
            return True

    @classmethod
    async def update_fee(cls, asin: str, fee_rate: float, shipping_fee: int):
        async with cls.session_scope() as session:
            stmt = update(cls).where(cls.asin == asin)\
                   .values(fee_rate=fee_rate, shipping_fee=shipping_fee)
            result = await session.execute(stmt)
            result.close()
            return True

    @classmethod
    async def delete_rows(cls, filename: str):
        async with cls.session_scope() as session:
            stmt = delete(cls).where(cls.filename == filename)
            await session.execute(stmt)
            return True

    @classmethod
    async def delete_rows_lower_price(cls, profit: int=200, profit_rate: float=0.1, unit_count: int=10, drops: int=3) -> True:
        async with cls.session_scope() as session:
            stmt = select(cls.asin, cls.filename).join(KeepaProducts, KeepaProducts.asin == cls.asin, isouter=True).where(or_(
                cls.profit < profit,
                cls.profit_rate < profit_rate,
                cls.unit > unit_count,
                KeepaProducts.sales_drops_90 <= drops,
            ))
            result = await session.execute(stmt)
            for asin, filename in result.all():
                stmt = delete(cls).where(cls.asin == asin, cls.filename == filename)
                await session.execute(stmt)
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


def init_db(engine):
    Base.metadata.create_all(bind=engine)
