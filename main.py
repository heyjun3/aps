import logging.config
import sys

from mws import api
from keepa import keepa
from crawler.buffalo import buffalo
from crawler.pc4u import pc4u
from crawler.rakuten import rakuten
from crawler.super import super
from settings import LOGGING_CONF_PATH


if __name__ == '__main__':
    logging.config.fileConfig(LOGGING_CONF_PATH, disable_existing_loggers=False)
    logger = logging.getLogger(__name__)

    args = sys.argv

    if args[1] == 'keepa':
        keepa.keepa_worker()
    elif args[1] == 'mws':
        api.main()
    elif args[1] == 'buffalo':
        buffalo.main()
    elif args[1] == 'pc4u':
        pc4u.schedule_pc4u_task_everyday()
    elif args[1] == 'rakuten':
        rakuten.schedule()
    elif args[1] == 'super':
        super.schedule_super_task()
    else:
        sys.stdout.write(f'{args[1]} is not a command')
