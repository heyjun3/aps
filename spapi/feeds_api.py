import urllib.parse
from functools import partial

import settings
import log_settings
from spapi.spapi import SPAPI

logger = log_settings.get_logger(__name__)


class FeedsAPI(SPAPI):

    def __init__(self) -> None:
        super().__init__()

    async def create_feed_document(self, content_type: str, encoding: str):
        return await self._request(partial(self._create_feed_document, content_type, encoding))

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
