import urllib.parse
from functools import partial
from io import StringIO
from typing import Callable

import requests
from requests import Response

import settings
import log_settings
from spapi.spapi import SPAPI

logger = log_settings.get_logger(__name__)


def logger_decorator_with_response(func: Callable) -> Callable:
    async def _logger_decorator(self, *args, **kwargs):
        logger.info({'action': func.__name__, 'status': 'run',
                    'args': args, 'kwargs': kwargs})
        result = await func(self, *args, **kwargs)
        logger.info({'action': func.__name__,
                    'status': 'done', 'response': result})
        return result
    return _logger_decorator


class FeedsAPI(SPAPI):

    def __init__(self) -> None:
        super().__init__()

    async def upload_feed(self, url: str, filename: str, file: StringIO, content_type: str) -> Response:
        access_token = await self.get_spapi_access_token()
        headers = self.create_authorization_headers(access_token, 'POST', url)
        res = requests.post(url, headers=headers, files={
                            'file': (filename, file, content_type)})
        return res

    @logger_decorator_with_response
    async def create_feed_document(self, content_type: str, encoding: str):
        return await self._request(partial(self._create_feed_document, content_type, encoding))

    @logger_decorator_with_response
    async def create_feed(self, feed_type: str, document_id: str) -> dict:
        return await self._request(partial(self._create_feed, feed_type, document_id))

    @logger_decorator_with_response
    async def get_feed(self, feed_id: str) -> dict:
        return await self._request(partial(self._get_feed, feed_id))

    @logger_decorator_with_response
    async def get_feed_document(self, document_id: str) -> dict:
        return await self._request(partial(self._get_feed_document, document_id))

    def _create_feed_document(self, content_type: str, encoding: str = 'UTF-8') -> dict:
        method = 'POST'
        path = '/feeds/2021-06-30/documents'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        body = {
            'contentType': f'{content_type}; charset={encoding}'
        }

        return (method, url, None, body)

    def _create_feed(self, feed_type: str, document_id: str) -> tuple:
        method = 'POST'
        path = '/feeds/2021-06-30/feeds'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        body = {
            'feedType': feed_type,
            'marketplaceIds': [
                self.marketplace_id,
            ],
            'inputFeedDocumentId': document_id,
        }
        return (method, url, None, body)

    def _get_feed(self, feed_id: str) -> tuple:
        method = 'GET'
        path = f'/feeds/2021-06-30/feeds/{feed_id}'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        return (method, url, None, None)

    def _get_feed_document(self, document_id: str) -> tuple:
        method = 'GET'
        path = f'/feeds/2021-06-30/documents/{document_id}'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        return (method, url, None, None)
