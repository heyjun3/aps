import sys
from multiprocessing import Process

from mws import api
from keepa import keepa
from spapi import spapi_tasks
from spapi.spapi_tasks import UpdatePriceAndRankTask
from crawler.buffalo import buffalo
from crawler.pc4u import pc4u
from crawler.rakuten import rakuten_tasks
from crawler.super import super_tasks
from crawler.netsea import netsea_tasks
from ims import repeat
from ims import monthly



if __name__ == '__main__':

    args = sys.argv

    if args[1] == 'keepa':
        keepa.main()
    elif args[1] == 'mws':
        process_get_price = Process(target=api.run_get_lowest_priced_offer_listtings_for_asin, daemon=True)
        process_get_fees = Process(target=spapi_tasks.run_get_my_fees_estimate_for_asin, daemon=True)
        process_get_price.start()
        process_get_fees.start()
        spapi_tasks.run_list_catalog_items()
        process_get_price.join()
        process_get_fees.join()
    elif args[1] == 'buffalo':
        buffalo.main()
    elif args[1] == 'pc4u':
        pc4u.main()
    elif args[1] == 'rakuten':
        rakuten_tasks.run_rakuten_search_all()
    elif args[1] == 'super':
        super_tasks.run_super_all_shops()
    elif args[1] == 'netsea':
        netsea_tasks.run_netsea_all_products()
    elif args[1] == 'repeat':
        repeat.main()
    elif args[1] == 'monthly':
        monthly.main()
    elif args[1] == 'spapi':
        UpdatePriceAndRankTask().main()
    else:
        sys.stdout.write(f'{args[1]} is not a command')
