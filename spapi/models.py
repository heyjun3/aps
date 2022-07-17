import datetime
from contextlib import contextmanager
from typing import List

from sqlalchemy import ForeignKey
from sqlalchemy import create_engine
from sqlalchemy import Column
from sqlalchemy import Float 
from sqlalchemy import String
from sqlalchemy import BigInteger
from sqlalchemy import Date
from sqlalchemy.orm import sessionmaker
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.dialects.postgresql import Insert
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

    # def upsert(self) -> True:
    #     with session_scope() as session:
    #         stmt = Insert(AsinsInfo).values(self.values)
    #         stmt = stmt.on_conflict_do_update(index_elements=['asin'], set_=self.values)
    #         session.execute(stmt)
    #     return True

    async def upsert(self) -> True:
        async with self.session_scope() as session:
            stmt = Insert(AsinsInfo).values(self.values)
            stmt = stmt.on_conflict_do_update(index_elements=['asin'], set_=self.values)
            await session.execute(stmt)
        return True

    # @classmethod 
    # def get(cls, jan: str, interval_days: int=30) -> List[dict]|None:
    #     date = (datetime.date.today() - datetime.timedelta(days=interval_days))
    #     with session_scope() as session:
    #         asins = session.query(cls).filter(cls.jan == jan, cls.modified > date).all()
    #         if asins:
    #             asins = [asin.values for asin in  asins]
    #             return asins
    #         else:
    #             return None

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
    def get_title(cls, asin: str) -> str|None:
        with session_scope() as session:
            title = session.query(cls.title).filter(cls.asin == asin).first()
            if title:
                return title[0]
            return None

    @classmethod
    async def get_title(cls, asin: str) -> str|None:
        async with cls.session_scope() as session:
            stmt = select(cls.title).where(cls.asin == asin)
            result = await session.execute(stmt)
            title = result.scalar()
        if title:
            return title
        return None

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

    def upsert(self) -> True:
        with session_scope() as session:
            stmt = Insert(SpapiPrices).values(self.values)
            stmt = stmt.on_conflict_do_update(index_elements=['asin'], set_=self.values)
            session.execute(stmt)
        return True

    async def upsert(self) -> True:
        async with self.session_scope() as session:
            stmt = Insert(SpapiPrices).values(self.values)
            stmt = stmt.on_conflict_do_update(index_elements=['asin'], set_=self.values)
            await session.execute(stmt)
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

    def upsert(self) -> True:
        with session_scope() as session:
            stmt = Insert(SpapiFees).values(self.values)
            stmt = stmt.on_conflict_do_update(index_elements=['asin'], set_=self.values)
            session.execute(stmt)
        return True

    async def upsert(self) -> True:
        async with self.session_scope() as session:
            stmt = Insert(SpapiFees).values(self.value)
            stmt = stmt.on_conflict_do_update(index_elements=['asin'], set_=self.values)
            await session.execute(stmt)
        return True

    # @classmethod
    # def get(cls, asin: str, interval_days: int=30) -> dict|None:
    #     date = datetime.date.today() - datetime.timedelta(days=interval_days)
    #     with session_scope() as session:
    #         asin_fee = session.query(cls).filter(cls.asin == asin, cls.modified > date).first()
    #         if asin_fee:
    #             return asin_fee.values
    #         else:
    #             return None

    @classmethod
    async def get(cls, asin: str, interval_days: int=30) -> dict|None:
        date = datetime.date.today() - datetime.timedelta(days=interval_days)
        async with cls.session_scope() as session:
            stmt = select(cls).where(cls.asin == asin, cls.modified > date)
            result = await session.execute(stmt)
            asin_fee = result.scalar()
            if asin_fee:
                return asin_fee.values
            return None

    @property
    def values(self):
        return {
            'asin': self.asin,
            'fee_rate': self.fee_rate,
            'ship_fee': self.ship_fee,
            'modified': self.modified,
        }


@contextmanager
def session_scope():
    try:
        session = Session()
        yield session
        session.commit()
    except Exception as ex:
        logger.error(ex)
        session.rollback()


def init_db():
    Base.metadata.create_all(bind=engine)
