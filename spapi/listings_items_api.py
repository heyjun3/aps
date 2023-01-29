import urllib.parse
from typing import List
from functools import partial

from spapi.spapi import SPAPI
import log_settings
import settings


logger = log_settings.get_logger(__name__)


class ListingsItemsAPI(SPAPI):

    def __init__(self) -> None:
        super().__init__()

    async def create_new_sku(self, *args, **kwargs):
        return await self._request(partial(self._patch_listings_item, *args, **kwargs))
    
    async def get_listing_item(self, *args):
        return await self._request(partial(self._get_listing_item, *args))

    def _get_listing_item(self, sku: str) -> dict:
        logger.info({'action': 'get_listing_item', 'status': 'run'})

        method = 'GET'
        path = f'/listings/2021-08-01/items/{self.seller_id}/{sku}'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            'marketplaceIds': self.marketplace_id,
            'issueLocale': 'ja_JP',
            'includedData': 'attributes'
        }

        logger.info({'action': 'get_listing_item', 'status': 'done'})
        return (method, url, query, None)

    def _patch_listings_item(self, sku: str, price: float, asin: str, 
                                  condition_note: str, product_type: str='PRODUCT', 
                                  condition_type: str='new_new') -> dict:
        logger.info({'action': 'patch_listings_item', 'status': 'run'})

        method = 'PATCH'
        path = f'/listings/2021-08-01/items/{self.seller_id}/{sku}'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            'marketplaceIds': self.marketplace_id,
            'issueLocale': 'en_US',
        }
        body = {
            'productType': product_type,
            'patches': [
                {
                    'op': 'add',
                    'path': '/attributes/purchasable_offer',
                    'value': [{
                        'marketplace_id': self.marketplace_id,
                        'currency': 'JPY',
                        'our_price': [{
                            'schedule': [{
                                'value_with_tax': price,
                            }]
                        }]
                    }]
                },
                {
                    'op': 'add',
                    'path': '/attributes/merchant_suggested_asin',
                    'value': [{
                        'value': asin,
                        'marketplace_id': self.marketplace_id,
                    }]
                },
                {
                    'op': 'add',
                    'path': '/attributes/condition_type',
                    'value': [{
                        'value': condition_type,
                        'marketplace_id': self.marketplace_id,
                    }]
                },
                {
                    'op': 'add',
                    'path': '/attributes/condition_note',
                    'value': [{
                        'language_tag': 'ja_JP',
                        'value': condition_note,
                        'marketplace_id': self.marketplace_id,
                    }]
                },
                {
                    'op': 'add',
                    'path': '/attributes/fulfillment_availability',
                    'value': [{
                        'fulfillment_channel_code': 'AMAZON_JP',
                        'marketplace_id': self.marketplace_id,
                    }]
                },
                {
                    'op': 'add',
                    'path': '/attributes/batteries_required',
                    'value': [{
                        'value': 'false',
                        'marketplace_id': self.marketplace_id,
                    }]
                },
                {
                    'op': 'add',
                    'path': '/attributes/supplier_declared_dg_hz_regulation',
                    'value': [{
                        'value': 'not_applicable',
                        'marketplace_id': self.marketplace_id,
                    }]
                },
            ],
        }

        logger.info({'action': 'patch_listing_items', 'action': 'done'})
        return (method, url, query, body)

    def _put_listings_item(self, sku: str, asin: str, price: float, condition_type: str='new_new', condition_note: str=None) -> dict:
        logger.info({'action': 'put_listings_item', 'status': 'run'})

        method = 'PUT'
        path = f'/listings/2021-08-01/items/{self.seller_id}/{sku}'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            'marketplaceIds': self.marketplace_id,
            'issueLocale': 'en_US',
        }
        body = {
            'productType': 'SCREEN_PROTECTOR',
            'attributes': {
                 'purchasable_offer': [{
                    'currency': 'JPY',
                    'our_price': [{
                        'schedule': [{
                            'value_with_tax': price,
                        }]
                    }],
                    'marketplace_id': self.marketplace_id,
                }],
                'merchant_suggested_asin': [{
                    'value': asin,
                    'marketplace_id': self.marketplace_id,
                }],
                'condition_type': [{
                    'value': condition_type,
                    'marketplace_id': self.marketplace_id,
                }],
                'condition_note': [{
                    'language_tag': 'ja_JP',
                    'value': condition_note,
                    'marketplace_id': self.marketplace_id,
                }],
                'fulfillment_availability': [{
                    'fulfillment_channel_code': 'AMAZON_JP',
                    'marketplace_id': self.marketplace_id,
                }],
                'batteries_required': [{
                    'value': 'false',
                    'marketplace_id': self.marketplace_id,
                }],
                'supplier_declared_dg_hz_regulation': [{
                    'value': 'not_applicable',
                    'marketplace_id': self.marketplace_id,
                }]
            }
        }

        logger.info({'action': 'put_listings_item', 'status': 'done'})
        return (method, url, query, body)

    def search_definitions_product_types(self) -> tuple:
        logger.info({'action': 'search_difinitions_product_types', 'status': 'run'})
        
        method = 'GET'
        path = '/definitions/2020-09-01/productTypes'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            # 'keywords': ','.join(keywords),
            'marketplaceIds': self.marketplace_id,
        }

        logger.info({'action': 'search_difinitions_product_types', 'status': 'done'})
        return (method, url, query, None)

    async def get_definitinos_product_type(self, product_type: str) -> dict:
        logger.info({'action': 'get_difinitions_product_type', 'status': 'run'})

        method = 'GET'
        path = f'/definitions/2020-09-01/productTypes/TELEVISION'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            'sellerId': self.seller_id,
            'marketplaceIds': self.marketplace_id,
        }

        response = await self.request(method, url, query)

        logger.info({'action': 'get_difinitions_product_type', 'status': 'done'})
        return response