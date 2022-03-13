import datetime
from contextlib import contextmanager

from sqlalchemy import create_engine, null
from sqlalchemy import Column
from sqlalchemy import String
from sqlalchemy import BigInteger
from sqlalchemy import Date
from sqlalchemy.orm import sessionmaker
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.dialects.postgresql import Insert

import log_settings
import settings

engine = create_engine(settings.DB_URL)
Session = sessionmaker(bind=engine)
Base = declarative_base()
logger = log_settings.get_logger(__name__)


class AsinsInfo(Base):
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

    def save(self):
        with session_scope() as session:
            session.add(self)
        return True

    def upsert(self):
        with session_scope() as session:
            stmt = Insert(AsinsInfo).values(self.values)
            stmt = stmt.on_conflict_do_update(index_elements=['asin'], set_=self.values)
            session.execute(stmt)
        return True

    @classmethod 
    def get(cls, jan: str, interval_days: int=30) -> list|None:
        date = (datetime.date.today() - datetime.timedelta(days=interval_days))
        with session_scope() as session:
            asins = session.query(cls).filter(cls.jan == jan, cls.modified < date).all()
            if asins:
                return asins
            else:
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