from contextlib import contextmanager
from logging import getLogger

from sqlalchemy import create_engine
from sqlalchemy import Column
from sqlalchemy import String
from sqlalchemy import Integer
from sqlalchemy import Float
from sqlalchemy.orm import sessionmaker
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.declarative import declarative_base

import settings

engine = create_engine(settings.DB_URL)
Session = sessionmaker(bind=engine)
Base = declarative_base()
logger = getLogger(__name__)


class MWS(Base):
    __tablename__ = 'mws_products'
    asin = Column(String, primary_key=True, nullable=False)
    filename = Column(String, primary_key=True, nullable=False)
    title = Column(String)
    jan = Column(String)
    unit = Column(Integer)
    price = Column(Integer)
    fee_rate = Column(Float)
    shipping_fee = Column(Integer)
    profit = Column(Integer)
    profit_rate = Column(Float)

    def save(self):
        with session_scope() as session:
            try:
                session.add(self)
            except IntegrityError as ex:
                logger.error(ex)
                return False
            return True

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
    def update_fee_and_profit(cls, asin: str, filename: str, fee_rate: float, shipping_fee: int, profit: int, profit_rate: float):
        with session_scope() as session:
            mws = session.query(cls).filter(cls.asin == asin, cls.filename == filename).first()
            try:
                mws.fee_rate = fee_rate
                mws.shipping_fee = shipping_fee
                mws.profit = profit
                mws.profit_rate = profit_rate
            except Exception as ex:
                logger.error(f'action=update_fee_and_profit error={ex}')
                return False
            return True

    @property
    def value(self):
        return {
            'asin': self.asin,
            'jan': self.jan,
            'unit': self.unit,
            'price': self.price,
            'fee_rate': self.fee_rate,
            'shipping_fee': self.shipping_fee,
        }


@contextmanager
def session_scope():
    session = Session()
    session.expire_on_commit = False
    try:
        yield session
        session.commit()
    except Exception as ex:
        logger.error(f'action=session_scope error={ex}')
        session.rollback()
    finally:
        session.expire_on_commit = True


def init_db():
    Base.metadata.create_all(bind=engine)
