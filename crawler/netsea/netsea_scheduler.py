import time
import multiprocessing
from typing import Callable

import schedule

from crawler.netsea import netsea_tasks


def run_process(job_func: Callable) -> None:
    process = multiprocessing.Process(target=job_func)
    process.start()
    process.join()


def main() -> None:

    schedule.every().day.at('01:00').do(run_process, netsea_tasks.run_get_discount_products)

    while True:
        schedule.run_pending()
        time.sleep(1)
