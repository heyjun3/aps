import argparse
import sys

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
from ims import repeat
from ims import monthly


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('task', help='run task you use here', type=str)
    parser.add_argument('-i', '--id', help='Enter shop id', type=str, default=None)
    args = parser.parse_args()
    task = args.task
    shop_id = args.id

    match (task, shop_id):
        case ('keepa', None):
            keepa.main()
        case ('amz', None):
            RunAmzTask().main()
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
            UpdatePriceAndRankTask().main()
        case ('pcones', None):
            pcones.main()
        case ('mws', None):
            MWS.delete_rows_lower_price()
        case _:
            sys.stdout.write(f'{task} is not a command')
