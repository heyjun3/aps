import time

import requests
from requests import Session, Response

import log_settings


logger = log_settings.get_logger(__name__)


HEADERS = {"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36Mozilla/5.0 (Windows NT 10.0; Win64; x64) "}


class RequestException(Exception):
    pass


def request(url: str, method: str = 'GET', session: Session = requests.Session(), data: dict = None) -> Response:
    for _ in range(60):
        try:
            response = session.request(method=method, url=url, timeout=60.0, headers=HEADERS, data=data)
            if not response.status_code == 200:
                logger.error(response.status_code)
                raise RequestException
            return response
        except Exception as e:
            logger.error(f'action=request error={e}')
            logger.error(url)
            time.sleep(30)
