from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import sessionmaker

import settings


class ModelsBase(object):

    engine = create_async_engine(settings.DB_ASYNC_URL, future=True)
    async_session = sessionmaker(engine, expire_on_commit=False, class_=AsyncSession)
