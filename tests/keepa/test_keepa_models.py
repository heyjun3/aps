import pytest
import pytest_asyncio
from sqlalchemy.ext.asyncio import create_async_engine
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import sessionmaker

import settings
from keepa.models import Base, KeepaProducts
from keepa.models import KeepaProducts


class TestModels(object):

    @pytest_asyncio.fixture(autouse=True)
    async def fixture(cls):
        engine = create_async_engine(settings.DB_TEST_URL)
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.create_all)
        KeepaProducts.host_url = settings.DB_TEST_URL
        await KeepaProducts(asin='test', sales_drops_90=10, price_data={1000: 1000}, rank_data={2000: 2000}).save()
        yield
        async with engine.begin() as conn:
            await conn.run_sync(Base.metadata.drop_all)

    @pytest.mark.asyncio
    async def test_get_async(self):
        await KeepaProducts(asin='test1', sales_drops_90=10, price_data={1000: 1000}, rank_data={2000: 2000}).save()
        result = await KeepaProducts.get_keepa_products_by_asins(['test1'])
        assert result[0].asin == 'test1'
        assert result[0].sales_drops_90 == 10
        assert result[0].price_data == {'1000': 1000}
        assert result[0].rank_data == {'2000': 2000}

    @pytest.mark.asyncio
    async def test_insert_all(self):
        products = [
            KeepaProducts(asin='asin', price_data={'1111': 1111}, rank_data={'2222': 2222}),
            KeepaProducts(asin='test', price_data={'1111': 1111}, rank_data={'2222': 2222}, render_data={'3333': 3333}),
        ]
        result = await KeepaProducts.insert_all_on_conflict_do_update_chart_data(products)
        assert result == True
        result = await KeepaProducts.get_keepa_products_by_asins(['test'])
        assert result[0].asin == 'test'
        assert result[0].price_data == {'1111': 1111}
        assert result[0].rank_data == {'2222': 2222}
        assert result[0].render_data == {'3333': 3333}
    