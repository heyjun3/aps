from __future__ import annotations
import datetime
from typing import List

from sqlalchemy import ForeignKey
from sqlalchemy import create_engine
from sqlalchemy import Column
from sqlalchemy import Float 
from sqlalchemy import String
from sqlalchemy import BigInteger
from sqlalchemy import Date
from sqlalchemy import bindparam
from sqlalchemy import update
from sqlalchemy.orm import sessionmaker
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.dialects.postgresql import Insert
from sqlalchemy.dialects.postgresql import insert
from sqlalchemy.future import select

import log_settings
import settings
from models_base import ModelsBase

engine = create_engine(settings.DB_URL)
Session = sessionmaker(bind=engine)
Base = declarative_base()
logger = log_settings.get_logger(__name__)


class AsinsInfo(Base, ModelsBase):
    __tablename__= 'asins_info'
    asin = Column(String, primary_key=True, nullable=False)
    jan = Column(String)
    title = Column(String)
    quantity = Column(BigInteger)
    modified = Column(Date, onupdate=datetime.date.today(), default=datetime.date.today())

    def __init__(self, asin: str=None, jan: str=None, title: str=None, quantity: int=None):
        self.modified = datetime.date.today()
        self.asin = asin
        self.jan = jan
        self.title = title
        self.quantity = quantity

    def __repr__(self):
        return (f"{self.__class__.__name__}({self.asin}, {self.jan}, "
        f"{self.title}, {self.quantity}, {self.modified})")


    async def save(self) -> True:
        async with self.session_scope() as session:
            session.add(self)
        return True

    async def upsert(self) -> True:
        async with self.session_scope() as session:
            stmt = Insert(AsinsInfo).values(self.values)
            stmt = stmt.on_conflict_do_update(index_elements=['asin'], set_=self.values)
            await session.execute(stmt)
        return True

    @classmethod
    async def insert_all_on_conflict_do_update(cls, values: List[dict]) -> True|None:
        if not values:
            return
        stmt = insert(cls).values(values)
        do_update_stmt = stmt.on_conflict_do_update(
            index_elements=['asin'],
            set_=dict(
                jan=stmt.excluded.jan,
                title=stmt.excluded.title,
                quantity=stmt.excluded.quantity,
            ))
        async with cls.session_scope() as session:
            await session.execute(do_update_stmt)
            return True

    @classmethod
    async def get(cls, jan: str, interval_days: int=30) -> List[dict]|None:
        date = (datetime.date.today() - datetime.timedelta(days=interval_days))
        async with cls.session_scope() as session:
            stmt = select(cls).where(cls.jan == jan, cls.modified > date)
            result = await session.execute(stmt)
            asins = result.scalars().all()
        if asins:
            return [asin.values for asin in asins]
        return None

    @classmethod
    async def get_asin_object_by_jan_list(cls, jan_list: List[str]) -> List[AsinsInfo]:
        async with cls.session_scope() as session:
            stmt = select(cls).where(cls.jan.in_(jan_list))
            result = await session.execute(stmt)
            return result.scalars().all()

    @classmethod
    async def get_title(cls, asin: str) -> str|None:
        async with cls.session_scope() as session:
            stmt = select(cls.title).where(cls.asin == asin)
            result = await session.execute(stmt)
            title = result.scalar()
        if title:
            return title
        return None

    @classmethod
    async def get_asins_all(cls) -> List[str]|None:
        async with cls.session_scope() as session:
            stmt = select(cls.asin)
            result = await session.execute(stmt)
            asins = result.scalars().all()
        if asins:
            return asins
        return None

    @classmethod
    async def get_jan_code_all(cls) -> List[str]:
        async with cls.session_scope() as session:
            stmt = select(cls.jan)
            result = await session.execute(stmt)
            return result.scalars().all()

    @property
    def values(self):
        return {
            'asin': self.asin,
            'jan': self.jan,
            'title': self.title,
            'quantity': self.quantity,
            'modified': self.modified,
        }


class SpapiPrices(Base, ModelsBase):

    __tablename__ = 'spapi_prices'
    asin = Column(String, ForeignKey('asins_info.asin'), primary_key=True, nullable=False)
    price = Column(BigInteger)
    modified = Column(Date, default=datetime.date.today(), onupdate=datetime.date.today())

    def __init__(self, asin: str, price: int):
        self.asin = asin
        self.price = price
        self.modified = datetime.date.today()

    async def upsert(self) -> True:
        async with self.session_scope() as session:
            stmt = Insert(SpapiPrices).values(self.values)
            stmt = stmt.on_conflict_do_update(index_elements=['asin'], set_=self.values)
            await session.execute(stmt)
        return True

    @classmethod
    async def insert_all_on_conflict_do_update_price(cls, values: List[dict]) -> True:
        stmt = insert(cls).values([{
            'asin': value['asin'],
            'price': value['price'],
        } for value in values])
        update_on_stmt = stmt.on_conflict_do_update(
            index_elements=['asin'],
            set_=dict(
                price=stmt.excluded.price,
            )
        )
        async with cls.session_scope() as session:
            await session.execute(update_on_stmt)
            return True

    @property
    def values(self):
        return {
            'asin': self.asin,
            'price': self.price,
            'modified': self.modified,
        }


class SpapiFees(Base, ModelsBase):

    __tablename__ = 'spapi_fees'
    asin = Column(String, ForeignKey('asins_info.asin'), primary_key=True, nullable=False)
    fee_rate = Column(Float)
    ship_fee = Column(BigInteger)
    modified = Column(Date, default=datetime.date.today(), onupdate=datetime.date.today())

    def __init__(self, asin: str, fee_rate: float, ship_fee: int):
        self.asin = asin
        self.fee_rate = fee_rate
        self.ship_fee = ship_fee
        self.modified = datetime.date.today()

    async def upsert(self) -> True:
        async with self.session_scope() as session:
            stmt = Insert(SpapiFees).values(self.values)
            stmt = stmt.on_conflict_do_update(index_elements=['asin'], set_=self.values)
            await session.execute(stmt)
        return True

    @classmethod
    async def insert_all_on_conflict_do_update_fee(cls, values: List[SpapiFees|dict]) -> True|None:
        if not values:
            return
        stmt = insert(cls).values([{
            'asin': value.asin,
            'fee_rate': value.fee_rate,
            'ship_fee': value.ship_fee,
        } if isinstance(value, SpapiFees) else {
            "asin": value.get("asin"),
            "fee_rate": value.get("fee_rate"),
            "ship_fee": value.get("ship_fee"),
        } for value in values])
        
        update_on_stmt = stmt.on_conflict_do_update(
            index_elements=['asin'],
            set_=dict(
                fee_rate=stmt.excluded.fee_rate,
                ship_fee=stmt.excluded.ship_fee,
            )
        )
        async with cls.session_scope() as session:
            await session.execute(update_on_stmt)
            return True

    @classmethod
    async def get(cls, asin: str, interval_days: int=30) -> dict|None:
        date = datetime.date.today() - datetime.timedelta(days=interval_days)
        async with cls.session_scope() as session:
            stmt = select(cls).where(cls.asin == asin, cls.modified > date)
            result = await session.execute(stmt)
            asin_fee = result.scalar()
            if asin_fee:
                return asin_fee.values
            return 

    @classmethod
    async def get_asins_fee(cls, asins: List[str], interval_days: int=30) -> List[SpapiFees]:
        date = datetime.date.today() - datetime.timedelta(days=interval_days)
        async with cls.session_scope() as session:
            stmt = select(cls).where(cls.asin.in_(asins), cls.modified > date).order_by(cls.asin)
            result = await session.execute(stmt)
            return result.scalars().all()

    @classmethod
    async def get_asins_after_update_interval_days(cls, interval_days: int=30) -> List[str]:
        past_date = datetime.date.today() - datetime.timedelta(days=interval_days)
        async with cls.session_scope() as session:
            stmt = select(cls.asin).where(cls.modified > past_date)
            result = await session.execute(stmt)
            asins = result.scalars().all()
        return asins

    @classmethod
    async def get_asins_all(cls) -> List[str]:
        async with cls.session_scope() as session:
            stmt = select(cls.asin)
            result = await session.execute(stmt)
            asins = result.scalars().all()
            return asins

    @property
    def values(self):
        return {
            'asin': self.asin,
            'fee_rate': self.fee_rate,
            'ship_fee': self.ship_fee,
            'modified': self.modified,
        }


def init_db():
    Base.metadata.create_all(bind=engine)
