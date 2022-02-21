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


def register():
    schedule.every().day.at('09:00').do(run_threaded, netsea_tasks.run_new_product_search)
    schedule.every().day.at('05:00').do(run_threaded, super_tasks.run_schedule_super_task)
    schedule.every().day.at('17:00').do(run_threaded, pc4u.main)
    schedule.every().day.at('17:30').do(run_threaded, buffalo.main)
    # schedule.every().monday.at('01:00').do(run_threaded, netsea.new_shop_search)
    schedule.every().saturday.at('07:00').do(run_threaded, repeat.main)


if __name__ == '__main__':

    register()

    while True:
        schedule.run_pending()
        time.sleep(1)
