import argparse
import sys

from keepa import keepa
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
    args = parser.parse_args()
    task = args.task

    if task == 'keepa':
        keepa.main()
    elif task == 'mws':
        RunAmzTask().main()
    elif task == 'buffalo':
        buffalo.main()
    elif task == 'pc4u':
        pc4u.main()
    elif task == 'rakuten':
        rakuten_tasks.run_rakuten_search_all()
    elif task == 'super':
        super_tasks.run_super_all_shops()
    elif task == 'netsea':
        netsea_tasks.run_netsea_all_products()
    elif task == 'repeat':
        repeat.main()
    elif task == 'monthly':
        monthly.main()
    elif task == 'spapi':
        UpdatePriceAndRankTask().main()
    elif task == 'pcones':
        pcones.main()
    else:
        sys.stdout.write(f'{task} is not a command')
