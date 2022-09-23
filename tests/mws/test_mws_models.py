import pytest
from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import sessionmaker
import pytest_asyncio

from mws.models import MWS
from mws.models import Base
import settings


class TestModels(object):

    @pytest_asyncio.fixture(autouse=True)
    async def fixture(cls):
        engine = create_async_engine(settings.DB_TEST_URL)
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.create_all)
        MWS.host_url = settings.DB_TEST_URL
        await MWS(asin="TEST", filename='testfilename', title='testtitle', jan='testjan',
             unit=10, price=10000, cost=1000, fee_rate=0.1, shipping_fee=500).save()
        await MWS(asin='testprice', filename='testfileprice').save()
        yield
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.drop_all)

    @pytest.mark.asyncio
    async def test_save(self):
        result = await MWS(asin='TEST01', filename='TESTfilename').save()
        assert result == True

    @pytest.mark.asyncio
    async def test_get_filenames(self):
        result = await MWS.get_filenames()
        assert result == ['testfilename', 'testfileprice']

    @pytest.mark.asyncio
    async def test_get_price_is_None_products(self):
        result = await MWS.get_price_is_None_asins()
        assert result == ['testprice']

    @pytest.mark.asyncio
    async def test_get_fee_is_None_asins(self):
        result = await MWS.get_fee_is_None_asins()
        assert result == ['testprice']

    @pytest.mark.asyncio
    async def test_update_price(self):
        result = await MWS.update_price('testprice', 9999)
        assert result == True
        mws = await MWS.get('testprice')
        assert mws.price == 9999

    @pytest.mark.asyncio
    async def test_update_fee(self):
        result = await MWS.update_fee('testprice', 0.5, 1000)
        assert result == True
        mws = await MWS.get('testprice')
        assert mws.fee_rate == 0.5
        assert mws.shipping_fee == 1000

    @pytest.mark.asyncio
    async def test_delete_rows(self):
        result = await MWS.delete_rows('testfilename')
        assert result == True
        filenames = await MWS.get_filenames()
        assert filenames == ['testfileprice']
