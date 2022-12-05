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
    async def test_insert_all(self):
        records = [MWS(asin=f'test{i}', filename=f'testfilename{i}') for i in range(10)]
        result = await MWS.insert_all_on_conflict_do_nothing(records)
        assert result == True

    @pytest.mark.asyncio
    async def test_insert_all_update(self):
        records = [
            MWS(asin='TEST', filename='testfilename', price=2000),
            MWS(asin='aaaaa', filename='testfilename', price=3000),
        ]
        result = await MWS.insert_all_on_conflict_do_update_price(records)
        assert result == True
        mws = await MWS.get('TEST')
        assert mws.price == 2000
        mws = await MWS.get('aaaaa')
        assert mws.price == 3000

    @pytest.mark.asyncio
    async def test_insert_all_update_fee(self):
        records = [
            MWS(asin='test', filename='test', fee_rate=0.1, shipping_fee=1000),
            MWS(asin='test1', filename='test', fee_rate=1.1, shipping_fee=9000),
        ]
        result = await MWS.insert_all_on_conflict_do_update_fee(records)
        assert result == True
        mws = await MWS.get('test')
        assert mws.fee_rate == 0.1
        assert mws.shipping_fee == 1000

    @pytest.mark.asyncio
    async def test_bulk_update_prices(self):
        records = [
            {"asin": "TEST", "price": 9999},
            {"asin": "testprice", "price": 2222},
        ]
        result = await MWS.bulk_update_prices(records)
        assert result == True
        mws = await MWS.get('TEST')
        assert mws.price == 9999
        mws = await MWS.get('testprice')
        assert mws.price == 2222

    @pytest.mark.asyncio
    async def test_bulk_update_prices_failed(self):
        result = await MWS.bulk_update_prices([])
        assert result == None

    @pytest.mark.asyncio
    async def test_bulk_update_fees(self):
        records = [
            {"asin": "TEST", "fee_rate": 0.5, "ship_fee": 1000},
            {"asin": "testprice", "fee_rate": 0.1, "ship_fee": 2222},
        ]
        result = await MWS.bulk_update_fees(records)
        assert result == True
        mws = await MWS.get("TEST")
        assert mws.fee_rate == 0.5
        assert mws.shipping_fee == 1000
        mws = await MWS.get("testprice")
        assert mws.fee_rate == 0.1
        assert mws.shipping_fee == 2222

    @pytest.mark.asyncio
    async def test_bulk_update_failed(self):
        result = await MWS.bulk_update_fees([])
        assert result == None

    @pytest.mark.asyncio
    async def test_get_filenames(self):
        result = await MWS.get_filenames()
        assert result == ['testfilename']

    @pytest.mark.asyncio
    async def test_get_price_is_None_products(self):
        result = await MWS.get_object_by_price_is_None()
        assert result[0].asin == 'testprice'
        assert result[0].filename == 'testfileprice'

    @pytest.mark.asyncio
    async def test_get_asins_by_price_is_None(self):
        result = await MWS.get_asins_by_price_is_None()
        assert len(result) == 1
        assert result == ["testprice"]

    @pytest.mark.asyncio
    async def test_get_fee_is_None_asins(self):
        result = await MWS.get_fee_is_None_asins()
        assert result[0].asin == 'testprice'
        assert result[0].filename == 'testfileprice'

    @pytest.mark.asyncio
    async def test_get_asins_by_fee_is_None(self):
        result = await MWS.get_asins_by_fee_is_None()
        assert result == ["testprice",]

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
        assert filenames == []
