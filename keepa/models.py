from __future__ import annotations
import datetime
import itertools
from typing import List

from sqlalchemy import create_engine
from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy import func
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import Column, Integer, String, Date, or_
from sqlalchemy import JSON
from sqlalchemy.orm import sessionmaker
from sqlalchemy.orm.attributes import flag_modified
from sqlalchemy.exc import IntegrityError
from sqlalchemy.future import select
from sqlalchemy.dialects.postgresql import insert
from contextlib import contextmanager

import settings
import log_settings
from models_base import ModelsBase
from keepa import convert


logger = log_settings.get_logger(__name__)
logger_decorator = log_settings.decorator_logging(logger)
postgresql_engine = create_engine(settings.DB_URL, pool_size=20, max_overflow=0, pool_pre_ping=True)
async_engine = create_async_engine(settings.DB_ASYNC_URL, pool_timeout=3000)
Base = declarative_base()
NetseaBase = declarative_base()
Session = sessionmaker(bind=postgresql_engine, autoflush=True, expire_on_commit=False)


class KeepaProducts(Base, ModelsBase):
    __tablename__ = 'keepa_products'
    asin = Column(String, primary_key=True)
    sales_drops_90 = Column(Integer)
    created = Column(Date, default=datetime.date.today)
    modified = Column(Date, default=datetime.date.today, onupdate=datetime.date.today)
    price_data = Column(JSON)
    rank_data = Column(JSON)
    render_data = Column(JSON, default=convert.recharts_data, onupdate=convert.recharts_data)

    # [TODO]問題なければ消す

    # @classmethod
    # def object_get_db_asin(cls, asin, delay=30):
    #     with _session_scope() as session:
    #         delay_date = datetime.date.today() - datetime.timedelta(days=delay)
    #         product = session.query(cls).filter(cls.asin == asin, cls.modified >= delay_date).first()
    #         if product is None:
    #             return None
    #         return product

    @classmethod
    @logger_decorator
    async def get_keepa_products_by_asins(cls, asins: List[str]) -> List[KeepaProducts]|None:
        if not asins:
            return
        stmt = select(cls).where(cls.asin.in_(asins))
        async with cls.session_scope() as session:
            result = await session.execute(stmt)
            return result.scalars().all()

    @classmethod
    async def get_modified_count_by_date(cls, date: datetime.date) -> dict[str, int]:
        stmt = select(func.count(or_(cls.modified == date, None)), func.count(cls.modified))
        async with cls.session_scope() as session:
            result = await session.execute(stmt)
            modified, total = result.first()
            return {"modified": modified, "total": total}

    async def save(self):
        async with self.session_scope() as session:
            session.add(self)
            await session.commit()
        return True

    @classmethod
    def update_or_insert(cls, asin, drops, price_data, rank_data):
        with _session_scope() as session:
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

    # @classmethod
    # def update_price_and_rank_data(cls, asin:str, unix_time: float, price: int, rank: int) -> bool:
    #     logger.info('action=update_price_and_rank_data status=run')
    #     time = convert.unix_time_to_keepa_time(unix_time)
    #     with _session_scope() as session:
    #         product = session.query(cls).filter(cls.asin == asin).first()
    #         if product is None:
    #             return False

    #         product.price_data[time] = price
    #         product.rank_data[time] = rank
    #         flag_modified(product, 'price_data')
    #         flag_modified(product, 'rank_data')

    #     logger.info('action=update_price_and_rank_data status=done')
    #     return True

    @classmethod
    async def insert_all_on_conflict_do_update_chart_data(cls, products: List[KeepaProducts]) -> True|None:
        if not products:
            return

        stmt = insert(cls).values([{
            'asin': value.asin,
            'price_data': value.price_data,
            'modified': datetime.date.today(),
            'rank_data': value.rank_data,
            'render_data': value.render_data,
        } for value in products])
        update_do_stmt = stmt.on_conflict_do_update(
            index_elements=['asin'],
            set_=dict(
                price_data=stmt.excluded.price_data,
                modified=stmt.excluded.modified,
                rank_data=stmt.excluded.rank_data,
                render_data=stmt.excluded.render_data,
            )
        )
        async with cls.session_scope() as session:
            await session.execute(update_do_stmt)
            return True

    # @classmethod
    # def _update_price_and_rank_data(cls, asin:str, unix_time: float, price: int, rank: int) -> bool:
    #     time = convert.unix_time_to_keepa_time(unix_time)
    #     with _session_scope() as session:
    #         product = session.query(cls).filter(cls.asin == asin).first()
    #         product.price_data[time] = price
    #         product.rank_data[time] = rank
    #         product.render_data = convert.recharts_data({'price_data': product.price_data, 'rank_data': product.rank_data})
    #         return True

    @classmethod
    def get_keepa_product(cls, asin: str):
        with _session_scope() as session:
            product = session.query(cls).filter(cls.asin == asin).first()
            if product:
                return product
            else:
                return None

    @classmethod
    def get_products_not_modified(cls, count: int=864000):
        today = datetime.date.today()
        with _session_scope() as session:
            products = session.query(cls.asin).filter(cls.modified != today, cls.price_data != None, cls.rank_data != None).limit(count).all()
            return [product[0] for product in products]
    
    # @classmethod
    # def update_render_data(cls, asin: str):
    #     with _session_scope() as session:
    #         product = session.query(cls).filter(cls.asin == asin).first()
    #         context = {'price_data': product.price_data, 'rank_data': product.rank_data}
    #         product.render_data = convert.recharts_data(context)

    # @classmethod
    # async def async_update_render_data(cls, asin: str):
    #     async with cls.session_scope() as session:
    #         stmt = select(cls).where(cls.asin == asin)
    #         result = await session.execute(stmt)
    #         product = result.scalars().first()
    #         context = {'price_data': product.price_data, 'rank_data': product.rank_data}
    #         product.render_data = convert.recharts_data(context)
    #         await session.commit()

    # @classmethod
    # def set_render_data_all(cls):
    #     with _session_scope() as session:
    #         asin_list = session.query(cls.asin).all()
    #         asin_list = list(itertools.chain.from_iterable(asin_list))

    #     for asin in asin_list:
    #         cls.async_update_render_data(asin)

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
def _session_scope():
    session = Session()
    session.expire_on_commit = False
    try:
        yield session
        session.commit()
    except Exception as e:
        logger.error(f'action=session_scope error={e}')
        session.rollback()
        raise
    finally:
        session.expire_on_commit = True


def init_db():
    Base.metadata.create_all(bind=postgresql_engine)
