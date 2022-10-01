import time
import multiprocessing
import asyncio
from typing import Callable
from typing import Coroutine
from functools import partial

import schedule

import settings
from crawler.netsea import netsea_tasks
from crawler.super import super_tasks
from crawler.pc4u import pc4u
from crawler.buffalo import buffalo
from crawler.pcones import pcones
from crawler.rakuten import rakuten_tasks
from crawler.spread_sheet.spread_sheet import SpreadSheetCrawler
from mws.models import MWS


def run_coroutine_job(coroutine: Coroutine) -> None:
    asyncio.run(coroutine)

def run_process(job_func: Callable) -> None:
    process = multiprocessing.Process(target=job_func)
    process.start()


def main() -> None:

    schedule.every(30).minutes.do(run_process, partial(run_coroutine_job, MWS.delete_rows_lower_price()))

    schedule.every().day.at('01:00').do(run_process, netsea_tasks.run_get_discount_products)
    schedule.every().day.at('02:00').do(run_process, super_tasks.run_discount_product_search)
    schedule.every().day.at('05:00').do(run_process, super_tasks.run_schedule_super_task)
    schedule.every().day.at('09:00').do(run_process, netsea_tasks.run_new_product_search)
    schedule.every().day.at('16:00').do(run_process, pcones.main)
    schedule.every().day.at('17:00').do(run_process, pc4u.main)
    schedule.every().day.at('17:30').do(run_process, buffalo.main)
    schedule.every().day.at('18:00').do(run_process, SpreadSheetCrawler(settings.CREDENTIAL_FILE_NAME, settings.SHEET_TITLE, settings.SHEET_NAME).start_crawler())

    schedule.every().saturday.at('04:00').do(run_process, rakuten_tasks.run_rakuten_search_all)

    while True:
        schedule.run_pending()
        time.sleep(1)


if __name__ == '__main__':
    main()
