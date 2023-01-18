import logging
import os

import settings
from fluent import handler


def get_logger(name: str, level=logging.INFO) -> logging.getLogger:
    logger = logging.getLogger(name)
    logger.setLevel(level)
    formatter = logging.Formatter('[%(asctime)-15s]%(levelname)s:%(name)s:%(message)s')
    
    streamhandler = logging.StreamHandler()
    streamhandler.setFormatter(formatter)
    streamhandler.setLevel(level)

    logger.addHandler(streamhandler)
    return logger

def decorator_logging(logger: logging.Logger):
    def _wrapper(func):
        def _inner_wrapper(*args, **kwargs):
            logger.info({'action': func.__name__, 'status': 'run'})
            result = func(*args, **kwargs)
            logger.info({'action': func.__name__, 'status': 'done'})
            return result
        return _inner_wrapper
    return _wrapper
    