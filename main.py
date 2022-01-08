import logging.config
import sys

from mws import api
from keepa import keepa

from settings import LOGGING_CONF_PATH


if __name__ == '__main__':
    logging.config.fileConfig(LOGGING_CONF_PATH, disable_existing_loggers=False)
    logger = logging.getLogger(__name__)

    args = sys.argv

    if args[1] == 'keepa':
        keepa.keepa_worker()
    elif args[1] == 'mws':
        api.main()
    else:
        sys.stdout.write(f'{args[1]} is not a command')
