import argparse
import sys
import asyncio

from keepa import keepa
from mws.models import MWS
from spapi.spapi_tasks import UpdateChartDataRequestTask, UpdateChartData
from spapi.spapi_tasks import RunAmzTask
from spapi.register_service import RegisterService
from crawler.buffalo import buffalo
from crawler.pc4u import pc4u
from crawler.rakuten import rakuten_tasks
from crawler.rakuten import rakuten_scheduler
from crawler.super import super_tasks
from crawler.netsea import netsea_tasks
from crawler.pcones import pcones
from crawler.spread_sheet.spread_sheet import SpreadSheetCrawler
from crawler.paypaymall.paypaymoll import YahooShopCrawler
from ims import repeat
from ims import monthly
import settings


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('task', help='run task you use here', type=str)
    parser.add_argument('-i', '--id', help='Enter shop id', type=str, default=None)
    args = parser.parse_args()
    task = args.task
    shop_id = args.id

    match (task, shop_id):
        case ('keepa', None):
            asyncio.run(keepa.main())
        case ('amz', None):
            asyncio.run(RunAmzTask().main())
        case ('amz', 'queue'):
            asyncio.run(RunAmzTask().get_queue())
        case ('amz', 'catalog_item'):
            asyncio.run(RunAmzTask().search_catalog_items_v20220401())
        case ('amz', 'price'):
            asyncio.run(RunAmzTask().get_item_offers_batch())
        case ('amz', 'price_v2'):
            asyncio.run(RunAmzTask().get_item_offer())
        case ('amz', 'fees'):
            asyncio.run(RunAmzTask().get_my_fees_estimate())
        case ('buffalo', None):
            buffalo.main()
        case ('pc4u', None):
            pc4u.main()
        case ('rakuten', 'all'):
            rakuten_tasks.run_rakuten_search_all()
        case ('rakuten', 'scheduler'):
            rakuten_scheduler.main()
        case ('rakuten', _):
            rakuten_tasks.run_rakuten_search_at_shop_code(shop_id)
        case ('super', 'all'):
            super_tasks.run_super_all_shops()
        case ('super', 'new'):
            super_tasks.run_schedule_super_task()
        case ('super', 'discount'):
            super_tasks.run_discount_product_search()
        case ('super', 'bookmark'):
            super_tasks.run_get_favorite_products()
        case ('netsea', 'all'):
            netsea_tasks.run_netsea_all_products()
        case ('netsea', 'new'):
            netsea_tasks.run_new_product_search()
        case ('netsea', 'discount'):
            netsea_tasks.run_get_discount_products()
        case ('netsea', 'bookmark'):
            netsea_tasks.run_get_favorite_products()
        case ('netsea', _):
            netsea_tasks.run_netsea_at_shop_id(shop_id)
        case ('repeat', None):
            repeat.main()
        case ('monthly', None):
            monthly.main()
        case ('spapi', 'db'):
            asyncio.run(UpdateChartData().main())
        case ('spapi', 'request'):
            asyncio.run(UpdateChartDataRequestTask().main())
        case ('pcones', None):
            pcones.main()
        case ('yahoo', shop_id):
            YahooShopCrawler().search_by_shop_id(settings.YAHOO_APP_ID, shop_id)
        case ('mws', None):
            asyncio.run(MWS.delete_rows_lower_price())
        case ('spread_sheet', None):
            SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, settings.SHEET_TITLE, settings.SHEET_NAME).start_crawler()
        case ('spapi', "register"):
            asyncio.run(RegisterService(settings.CREDENTIAL_FILE_NAME).start_register("仕入帳_2023", "ADD"))
        case ('spapi', "check"):
            asyncio.run(RegisterService(settings.CREDENTIAL_FILE_NAME).check_registerd("仕入帳_2023", "ADD", "DB"))
        case ('spapi', "point"):
            asyncio.run(RegisterService(settings.CREDENTIAL_FILE_NAME).register_points([["4571585642054-B-23980-20231008", 1]]))
        case _:
            sys.stdout.write(f'{task} is not a command')
