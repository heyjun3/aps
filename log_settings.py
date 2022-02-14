import logging
import os

import settings


def get_logger(name: str, level=logging.INFO) -> logging.getLogger:
    logger = logging.getLogger(name)
    logger.setLevel(level)
    formatter = logging.Formatter('[%(asctime)-15s]%(levelname)s:%(name)s:%(message)s')

    filehander = logging.FileHandler(os.path.join(settings.BASE_PATH, 'logs', f'{name}.log'))
    filehander.setFormatter(formatter)
    filehander.setLevel(level)

    streamhandler = logging.StreamHandler()
    streamhandler.setFormatter(formatter)
    streamhandler.setLevel(level)

    logger.addHandler(filehander)
    logger.addHandler(streamhandler)

    return logger
    