import time
import threading

import schedule

from crawler.pc4u import pc4u
from crawler.buffalo import buffalo
from crawler.super import super_tasks
from crawler.netsea import netsea_tasks
from ims import repeat


def run_threaded(func):
    thread = threading.Thread(target=func)
    thread.start()


    # schedule.every().day.at('09:00').do(run_threaded, netsea_tasks.run_new_product_search)
    # schedule.every().day.at('00:10').do(run_threaded, netsea_tasks.run_get_discount_products)
    # schedule.every().day.at('05:00').do(run_threaded, super_tasks.run_schedule_super_task)
    # schedule.every().day.at('02:00').do(run_threaded, super_tasks.run_discount_product_search)
    # schedule.every().day.at('17:00').do(run_threaded, pc4u.main)
    # schedule.every().day.at('17:30').do(run_threaded, buffalo.main)
    # schedule.every().saturday.at('07:00').do(run_threaded, repeat.main)
