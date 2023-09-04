import urllib.parse
from functools import partial


import settings
import log_settings
from spapi.spapi import SPAPI
from spapi.utils import async_logger

logger = log_settings.get_logger(__name__)


class FeedsAPI(SPAPI):

    def __init__(self) -> None:
        super().__init__()

    @async_logger(logger)
    async def create_feed_document(self, content_type: str, encoding: str):
        return await self._request(partial(self._create_feed_document, content_type, encoding))

    @async_logger(logger)
    async def create_feed(self, feed_type: str, document_id: str) -> dict:
        return await self._request(partial(self._create_feed, feed_type, document_id))

    @async_logger(logger)
    async def get_feed(self, feed_id: str) -> dict:
        return await self._request(partial(self._get_feed, feed_id))

    @async_logger(logger)
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
