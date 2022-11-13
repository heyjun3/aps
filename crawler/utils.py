import time

from requests_html import HTMLSession, HTMLResponse

import log_settings


logger = log_settings.get_logger(__name__)
HEADERS = {"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36Mozilla/5.0 (Windows NT 10.0; Win64; x64) "}


class RequestException(Exception):
    pass


def request(url: str,
            method: str = 'GET',
            session: HTMLSession = HTMLSession(),
            data: dict = None,
            params: dict = None,
            time_sleep: int=0) -> HTMLResponse:

    for i in range(60):
        try:
            response = session.request(method=method, url=url, timeout=30.0, headers=HEADERS, data=data, params=params, allow_redirects=True)
            if response.status_code in (200, 404):
                time.sleep(time_sleep)
                return response
            logger.error({'messages': "request error", "status code": response.status_code})
            raise RecursionError
        except Exception as ex:
            logger.error({'messages': ex, 'request url': url})
            time.sleep(t if (t := 2 * i + 1) < 10 else 10)
