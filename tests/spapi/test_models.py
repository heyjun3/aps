import pytest
from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import sessionmaker
import pytest_asyncio

from spapi.models import AsinsInfo
from spapi.models import SpapiFees
from spapi.models import SpapiPrices
from spapi.models import Base
import settings

class TestModels(object):

    @pytest_asyncio.fixture(autouse=True)
    async def fixture(cls):
        engine = create_async_engine(settings.DB_TEST_URL)
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.create_all)
        AsinsInfo.url = settings.DB_TEST_URL
        SpapiFees.url = settings.DB_TEST_URL
        SpapiPrices.url = settings.DB_TEST_URL
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

    @pytest.mark.asyncio
    async def test_get(cls):
        obj = await AsinsInfo.get('9999')
        assert obj[0]['asin'] == 'TEST'
        assert obj[0]['title'] == 'test row'
        assert obj[0]['quantity'] == 1
        assert obj[0]['jan'] == '9999'

    @pytest.mark.asyncio
    async def test_upsert_insert(cls):
        await AsinsInfo('TEST1', '0000', 'test upsert', 2).upsert()
        obj = await AsinsInfo.get('0000')
        assert obj[0]['asin'] == 'TEST1'
        assert obj[0]['title'] == 'test upsert'
        assert obj[0]['quantity'] == 2 
        assert obj[0]['jan'] == '0000'

    @pytest.mark.asyncio
    async def test_upsert_update(cls):
        await AsinsInfo('TEST2', '1111', 'test update', 5).upsert()
        await AsinsInfo('TEST2', '2222', 'test update', 10).upsert()
        obj = await AsinsInfo.get('2222')
        assert obj[0]['asin'] == 'TEST2'
        assert obj[0]['title'] == 'test update'
        assert obj[0]['quantity'] == 10 
        assert obj[0]['jan'] == '2222'

    @pytest.mark.asyncio
    async def test_spapi_prices_upsert(cls):
        result = await SpapiPrices('NNNN', 4000).upsert()
        assert result == True

    @pytest.mark.asyncio
    async def test_spapi_fees_upsert(cls):
        result = await SpapiFees('NNNN', 0.1, 500).upsert()
        assert result == True

    @pytest.mark.asyncio
    async def test_spapi_fees_get(cls):
        result = await SpapiFees('TEST', 0.1, 1000).upsert()
        assert result == True
        fee = await SpapiFees.get('TEST')
        assert fee['asin'] == 'TEST'
        assert fee['fee_rate'] == 0.1
        assert fee['ship_fee'] == 1000
        