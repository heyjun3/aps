import pytest
from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import sessionmaker
import pytest_asyncio

from spapi.models import AsinsInfo
from spapi.models import Base
import settings

class TestModels(object):

    @pytest_asyncio.fixture(autouse=True)
    async def fixture(cls):
        engine = create_async_engine(settings.DB_TEST_URL)
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.create_all)
        AsinsInfo.engine = engine
        AsinsInfo.async_session = sessionmaker(AsinsInfo.engine, expire_on_commit=False, class_=AsyncSession)
        await AsinsInfo('TEST', '9999', 'test row', 1).save()
        yield
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.drop_all)

    @pytest.mark.asyncio
    async def test_get_title(cls):
        title = await AsinsInfo.get_title('TEST')
        title_none = await AsinsInfo.get_title('TESTNONE')
        assert title == 'test row'
        assert title_none == None

    


