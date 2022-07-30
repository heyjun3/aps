from contextlib import asynccontextmanager

from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import sessionmaker
from sqlalchemy.exc import IntegrityError

import settings
from log_settings import get_logger


logger = get_logger(__name__)


class ModelsBase(object):

    url = settings.DB_ASYNC_URL

    @classmethod
    @asynccontextmanager
    async def session_scope(cls):
        engine = create_async_engine(cls.url, future=True)
        async_session = sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)
        session = async_session()
        try:
            yield session
            await session.commit()
        except IntegrityError as ex:
            logger.error(ex)
            await session.rollback()
        except Exception as ex:
            logger.error(ex)
            await session.rollback()
        finally:
            await session.close()
