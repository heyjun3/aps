import logging.config
import os

from ims import monthly
from mws import api
from keepa import keepa


if __name__ == '__main__':
    logging.config.fileConfig(os.path.join(os.path.dirname(__file__), 'logging.conf'), disable_existing_loggers=False)
    logger = logging.getLogger(__name__)

    keepa.keepa_worker()
