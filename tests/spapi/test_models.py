import pytest
from sqlalchemy.ext.asyncio import create_async_engine
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
        AsinsInfo.host_url = settings.DB_TEST_URL
        SpapiFees.host_url = settings.DB_TEST_URL
        SpapiPrices.host_url = settings.DB_TEST_URL
        await AsinsInfo('TEST', '9999', 'test row', 1).save()
        await AsinsInfo('test', '19999', 'test row', 1).save()
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
    async def test_insert_all_on_conflict_do_update(cls):
        values = [
            {'asin': 'test1', 'jan': '1111', 'quantity': 1, 'title': 'title1'},
            {'asin': 'test2', 'jan': '2222', 'quantity': 10, 'title': 'title2'},
        ]
        result = await AsinsInfo.insert_all_on_conflict_do_update(values)
        assert result == True
        obj = await AsinsInfo.get('2222')
        assert obj[0]['asin'] == 'test2'
        assert obj[0]['title'] == 'title2'
        assert obj[0]['quantity'] == 10
        assert obj[0]['jan'] == '2222'

    @pytest.mark.asyncio
    async def test_insert_all_on_conflict_do_update_update_value(cls):
        values = [
            {'asin': 'test1', 'jan': '1111', 'quantity': 1, 'title': 'title1'},
            {'asin': 'test2', 'jan': '2222', 'quantity': 10, 'title': 'title2'},
        ]
        result = await AsinsInfo.insert_all_on_conflict_do_update(values)
        assert result == True
        values = [
            {'asin': 'test2', 'jan': '3333', 'quantity': 30, 'title': 'title3'},
        ]
        result = await AsinsInfo.insert_all_on_conflict_do_update(values)
        obj = await AsinsInfo.get('3333')
        assert obj[0]['asin'] == 'test2'
        assert obj[0]['title'] == 'title3'
        assert obj[0]['quantity'] == 30
        assert obj[0]['jan'] == '3333'

    @pytest.mark.asyncio
    async def test_spapi_prices_upsert(cls):
        result = await SpapiPrices('NNNN', 4000).upsert()
        assert result == True

    @pytest.mark.asyncio
    async def test_insert_all_conflict(cls):
        values = [
            {'asin': 'TEST', 'price': 1000},
            {'asin': 'test', 'price': 21000},
        ]
        result = await SpapiPrices.insert_all_on_conflict_do_update_price(values)
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

    @pytest.mark.asyncio
    async def test_insert_all_on_conflict_do_update_fee(self):
        values = [
            SpapiFees('test', 0.5, 1000),
            SpapiFees('TEST', 2, 2000),
        ]
        result = await SpapiFees.insert_all_on_conflict_do_update_fee(values)
        assert result == True
        fee = await SpapiFees.get('test')
        assert fee['asin'] == 'test'
        assert fee['fee_rate'] == 0.5
        assert fee['ship_fee'] == 1000
        