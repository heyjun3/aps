import argparse
import sys
import asyncio

from keepa import keepa
from mws.models import MWS
from spapi.spapi_tasks import UpdatePriceAndRankTask
from spapi.spapi_tasks import RunAmzTask
from crawler.buffalo import buffalo
from crawler.pc4u import pc4u
from crawler.rakuten import rakuten_tasks
from crawler.super import super_tasks
from crawler.netsea import netsea_tasks
from crawler.pcones import pcones
from crawler.spread_sheet.spread_sheet import SpreadSheetCrawler
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
        case ('buffalo', None):
            buffalo.main()
        case ('pc4u', None):
            pc4u.main()
        case ('rakuten', 'all'):
            rakuten_tasks.run_rakuten_search_all()
        case ('rakuten', _):
            rakuten_tasks.run_rakuten_search_at_shop_code(shop_id)
        case ('super', 'all'):
            super_tasks.run_super_all_shops()
        case ('super', 'new'):
            super_tasks.run_schedule_super_task()
        case ('super', 'discount'):
            super_tasks.run_discount_product_search()
        case ('netsea', 'all'):
            netsea_tasks.run_netsea_all_products()
        case ('netsea', 'new'):
            netsea_tasks.run_new_product_search()
        case ('netsea', 'discount'):
            netsea_tasks.run_get_discount_products()
        case ('netsea', _):
            netsea_tasks.run_netsea_at_shop_id(shop_id)
        case ('repeat', None):
            repeat.main()
        case ('monthly', None):
            monthly.main()
        case ('spapi', None):
            asyncio.run(UpdatePriceAndRankTask().main())
        case ('pcones', None):
            pcones.main()
        case ('mws', None):
            asyncio.run(MWS.delete_rows_lower_price())
        case ('spread_sheet', None):
            SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, settings.SHEET_TITLE, settings.SHEET_NAME).start_crawler()
        case _:
            sys.stdout.write(f'{task} is not a command')
