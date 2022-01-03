import logging
import time

import requests
from requests import Session, Response

logger = logging.getLogger(__name__)
logger.setLevel('INFO')

HEADERS = {"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36Mozilla/5.0 (Windows NT 10.0; Win64; x64) "}


class RequestException(Exception):
    pass


def request(url: str, session=requests.Session()) -> Response:
    for _ in range(60):
        try:
            response = session.get(url, timeout=60.0, headers=HEADERS)
            if not response.status_code == 200:
                raise RequestException
            return response
        except RequestException as e:
            logger.error(f'action=request error={e}')
            time.sleep(30)
