import urllib.parse
from functools import partial
from io import StringIO

import requests
from requests import Response

import settings
import log_settings
from spapi.spapi import SPAPI

logger = log_settings.get_logger(__name__)


class FeedsAPI(SPAPI):

    def __init__(self) -> None:
        super().__init__()

    async def upload_feed(self, url: str, filename: str, file: StringIO, content_type: str) -> Response:
        access_token = await self.get_spapi_access_token()
        headers = self.create_authorization_headers(access_token, 'POST', url)
        res = requests.post(url, headers=headers, files={'file': (filename, file, content_type)})
        return res

    async def create_feed_document(self, content_type: str, encoding: str):
        return await self._request(partial(self._create_feed_document, content_type, encoding))
    
    async def create_feed(self, feed_type: str, document_id: str) -> dict:
        return await self._request(partial(self._create_feed, feed_type, document_id))
    
    async def get_feed(self, feed_id: str) -> dict:
        return await self._request(partial(self._get_feed, feed_id))

    def _create_feed_document(self, content_type: str, encoding: str = 'UTF-8') -> dict:
        logger.info({'action': '_create_feed_document', 'status': 'run'})

        method = 'POST'
        path = '/feeds/2021-06-30/documents'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        body = {
            'contentType': f'{content_type}; charset={encoding}'
        }

        logger.info({'action': '_create_feed_document', 'status': 'done'})
        return (method, url, None, body)
    
    def _create_feed(self, feed_type: str, document_id: str) -> tuple:
        logger.info({'action': '_create_feed', 'status': 'run'})

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

        logger.info({'action': '_create_feed', 'status': 'done'})
        return (method, url, None, body)
    
    def _get_feed(self, feed_id: str) -> tuple:
        logger.info({'action': '_get_feed', 'status': 'run'})

        method = 'GET'
        path = f'/feeds/2021-06-30/feeds/{feed_id}'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)

        logger.info({'action': '_get_feed', 'status': 'done'})
        return (method, url, None, None)
