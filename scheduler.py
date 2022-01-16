import time
import threading
import logging.config

import schedule

from crawler.pc4u import pc4u
from crawler.buffalo import buffalo
from crawler.super import super
from crawler.netsea import netsea
from ims import repeat
from settings import LOGGING_CONF_PATH


def run_threaded(func):
    thread = threading.Thread(target=func)
    thread.start()


def register():
    schedule.every().day.at('09:00').do(run_threaded, netsea.new_product_search)
    schedule.every().day.at('05:00').do(run_threaded, super.schedule_super_task)
    schedule.every().day.at('17:00').do(run_threaded, pc4u.schedule_pc4u_task_everyday)
    schedule.every().day.at('17:30').do(run_threaded, buffalo.main)
    schedule.every().monday.at('01:00').do(run_threaded, netsea.new_shop_search)
    # schedule.every().saturday.at('07:00').do(run_threaded, repeatedly.main)


if __name__ == '__main__':
    logging.config.fileConfig(LOGGING_CONF_PATH, disable_existing_loggers=False)
    logger = logging.getLogger(__name__)

    register()

    while True:
        schedule.run_pending()
        time.sleep(1)
