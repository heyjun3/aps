import logging.config
import os

from ims import monthly
from mws import api
from keepa import keepa

from settings import LOGGING_CONF_PATH


if __name__ == '__main__':
    logging.config.fileConfig(LOGGING_CONF_PATH, disable_existing_loggers=False)
    logger = logging.getLogger(__name__)

    keepa.keepa_worker()
    # api.main()
