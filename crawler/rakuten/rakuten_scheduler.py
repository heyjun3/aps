import time
import multiprocessing
from typing import Callable

import schedule

from crawler.rakuten import rakuten_tasks


def run_process(job_func: Callable) -> None:
    process = multiprocessing.Process(target=job_func)
    process.start()
    process.join()


def main() -> None:

    schedule.every().saturday.at('04:00').do(run_process, rakuten_tasks.run_rakuten_search_all)

    while True:
        schedule.run_pending()
        time.sleep(1)
