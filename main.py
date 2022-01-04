import logging.config
import os

from ims import monthly


if __name__ == '__main__':
    logging.config.fileConfig(os.path.join(os.path.dirname(__file__), 'logging.conf'))
    logger = logging.getLogger(__name__)

    monthly.main()
