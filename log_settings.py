import logging
import os

import settings
from fluent import handler


def get_logger(name: str, level=logging.INFO) -> logging.getLogger:
    logger = logging.getLogger(name)
    logger.setLevel(level)
    formatter = logging.Formatter('[%(asctime)-15s]%(levelname)s:%(name)s:%(message)s')
    
    fluent_handler = handler.FluentHandler('app', 'localhost', 9880)
    fluent_format = {
        'host': '%(hostname)s',
        'time': '%(asctime)-15s',
        'level': '%(levelname)s',
        'name': '%(name)s',
        'message': '%(message)s',
    }
    fluent_handler.setFormatter(handler.FluentRecordFormatter(fluent_format))
    fluent_handler.setLevel(level)

    # docker 移行時に削除する。
    filehander = logging.FileHandler(os.path.join(settings.BASE_PATH, 'logs', f'{name}.log'))
    filehander.setFormatter(formatter)
    filehander.setLevel(level)

    streamhandler = logging.StreamHandler()
    streamhandler.setFormatter(formatter)
    streamhandler.setLevel(level)

    logger.addHandler(streamhandler)
    logger.addHandler(fluent_handler)

    return logger
    