import datetime
import hmac
import hashlib
import json
import os
import time
import urllib.parse
import re
from typing import List
import asyncio

import redis
import requests
import aiohttp

import settings
import log_settings


logger = log_settings.get_logger(__name__)


def logger_decorator(func):
    def _logger_decorator(*args, **kwargs):
        logger.info({'action': func.__name__, 'status': 'run'})
        result = func(*args, **kwargs)
        logger.info({'action': func.__name__, 'status': 'done'})
        return result
    return _logger_decorator


redis_client = redis.Redis(
        host=settings.REDIS_HOST,
        port=settings.REDIS_PORT,
        db=settings.REDIS_DB,
    )
ENDPOINT = 'https://sellingpartnerapi-fe.amazon.com'


async def request(method: str, url: str, params: dict=None, headers: dict=None, body: dict=None) -> aiohttp.ClientResponse:
    for _ in range(60):
        async with aiohttp.request(method, url, params=params, headers=headers, json=body) as response:
            if response.status == 200 and response is not None:
                response = await response.json()
                return response
            else:
                logger.error(response)
                await asyncio.sleep(10)


class SPAPI:

    def __init__(self):
        self.refresh_toke = settings.REFRESH_TOKEN
        self.client_id = settings.CLIENT_ID
        self.client_secret = settings.CLIENT_SECRET
        self.aws_secret_key = settings.AWS_SECRET_KEY
        self.aws_access_key = settings.AWS_ACCESS_ID
        self.marketplace_id = settings.MARKETPLACEID

    def sign(self, key, msg):
        return hmac.new(key, msg.encode('utf-8'), hashlib.sha256).digest()

    def get_signature_key(self, key, dateStamp, regionName, serviceName):
        kDate = self.sign(('AWS4' + key).encode('utf-8'), dateStamp)
        kRegion = self.sign(kDate, regionName)
        kService = self.sign(kRegion, serviceName)
        kSigning = self.sign(kService, 'aws4_request')
        return kSigning

    async def get_spapi_access_token(self, timeout_sec: int = 3500):

        access_token = redis_client.get('access_token')
        if access_token is not None:
            return access_token.decode()

        URL = 'https://api.amazon.com/auth/o2/token'
        headers = {
        'Host': 'api.amazon.com',
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
        }
        params = {
        'grant_type': 'refresh_token',
        'refresh_token': self.refresh_toke,
        'client_id': self.client_id,
        'client_secret': self.client_secret,
        }
        response = await request('POST', URL, params, headers)

        access_token = response.get('access_token')
        
        if access_token is None:
            text = await response.text()
            logger.error(text)
            raise Exception

        redis_client.set('access_token', access_token, ex=timeout_sec)

        return access_token

    async def create_authorization_headers(self, method: str, url: str, params: dict={}, body: dict={}) -> dict:
        region = 'us-west-2'
        service = 'execute-api'
        algorithm = 'AWS4-HMAC-SHA256'
        signed_headers = 'host;user-agent;x-amz-access-token;x-amz-date'
        user_agent = 'My SPAPI Client tool /1.0(Language=python/3.10)'

        host = urllib.parse.urlparse(ENDPOINT).netloc
        canonical_uri = urllib.parse.urlparse(url).path

        if body:
            body = json.dumps(body)
        else:
            body = ''

        utcnow = datetime.datetime.utcnow()
        amz_date = utcnow.strftime('%Y%m%dT%H%M%SZ')
        datestamp = utcnow.strftime('%Y%m%d')
        amz_access_token = await self.get_spapi_access_token()
        canonical_header_values = [host, user_agent, amz_access_token, f'{amz_date}\n']
        
        canonical_headers = '\n'.join([f'{head}:{value}' for head, value in zip(signed_headers.split(';'), canonical_header_values)])
        payload_hash = hashlib.sha256(body.encode('utf-8')).hexdigest()

        if params:
            canonical_querystring = urllib.parse.urlencode(sorted(params.items(), key=lambda x: (x[0], x[1])))
        else:
            canonical_querystring = ''

        canonical_request = '\n'.join([method, canonical_uri, canonical_querystring, canonical_headers, signed_headers, payload_hash])
        credential_scope = os.path.join(datestamp, region, service, 'aws4_request')
        signing_key = self.get_signature_key(self.aws_secret_key, datestamp, region, service)
        string_to_sign = '\n'.join([algorithm, amz_date, credential_scope, hashlib.sha256(canonical_request.encode('utf-8')).hexdigest()])
        
        signature = hmac.new(signing_key, (string_to_sign).encode('utf-8'), hashlib.sha256).hexdigest()
        authorization_header = f'{algorithm} Credential={self.aws_access_key}/{credential_scope}, SignedHeaders={signed_headers}, Signature={signature}'
        
        headers = {
            'host': urllib.parse.urlparse(ENDPOINT).netloc,
            'user-agent': user_agent,
            'x-amz-date': amz_date,
            'Authorization': authorization_header,
            'x-amz-access-token': amz_access_token,
        }

        return headers
    
    async def request(self, method: str, url: str, params: dict=None, body: dict=None) -> dict:
        headers = await self.create_authorization_headers(method, url, params, body)
        response = await request(method, url, params=params, body=body, headers=headers)
        return response

    async def get_my_fees_estimate_for_asin(self, asin: str, price: int, is_fba: bool = True, currency_code: str = 'JPY') -> dict:
        logger.info('action=get_my_fees_estimate_for_asin status=run')

        method = 'POST'
        path = f'/products/fees/v0/items/{asin}/feesEstimate'
        url = urllib.parse.urljoin(ENDPOINT, path)
        body =  {
            'FeesEstimateRequest': {
                'Identifier': asin,
                'PriceToEstimateFees': {
                    'ListingPrice': {
                        'Amount': price,
                        'CurrencyCode': currency_code
                    },
                },
                'IsAmazonFulfilled': is_fba,
                'MarketplaceId': self.marketplace_id, 
            }
        }
        response = await self.request(method, url, body=body)

        return response

    async def get_pricing(self, asin_list: list, item_type: str='Asin') -> dict:
        method = 'GET'
        path = '/products/pricing/v0/price'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'Asins': ','.join(asin_list),
            'ItemType': item_type,
            'MarketplaceId': self.marketplace_id,
        }
        response = await self.request(method, url, params=query)

        return response

    async def get_competitive_pricing(self, asin_list: list, item_type: str='Asin') -> dict:
        logger.info('action=get_competitive_pricing status=run')

        method = 'GET'
        path = '/products/pricing/v0/competitivePrice'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'MarketplaceId': self.marketplace_id,
            'Asins': ','.join(asin_list),
            'ItemType': item_type,
        }
        response = await self.request(method, url, params=query)

        logger.info('action=get_competitive_pricing status=done')
        return response

    async def get_item_offers(self, asin: str, item_condition: str='New') -> dict:
        logger.info('action=get_item_offers status=run')

        method = 'GET'
        path = f'/products/pricing/v0/items/{asin}/offers'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'MarketplaceId': self.marketplace_id,
            'ItemCondition': item_condition,
        }
        response = await self.request(method, url, params=query)

        logger.info('action=get_item_offers status=done')
        return response

    async def search_catalog_items(self, jan_list: list) -> dict:
        logger.info('action=search_catalog_items status=run')

        method = 'GET'
        path = '/catalog/2020-12-01/items'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'keywords': ','.join(jan_list),
            'marketplaceIds': self.marketplace_id,
            'includedData': 'identifiers,images,productTypes,salesRanks,summaries,variations'
        }
        response = await self.request(method, url, params=query)

        logger.info('action=search_catalog_items status=done')
        return response

    async def list_catalog_items(self, jan: str) -> dict:
        logger.info('action=list_catalog_items status=run')

        method = 'GET'
        path = '/catalog/v0/items'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'MarketplaceId': self.marketplace_id,
            'JAN': jan,
        }
        response = await self.request(method, url, params=query)

        logger.info('action=list_catalog_items status=done')
        return response

    async def get_catalog_item(self, asin: str) -> dict:
        logger.info('action=get_catalog_item status=run')

        method = 'GET'
        path = f'/catalog/2020-12-01/items/{asin}'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'marketplaceIds': self.marketplace_id,
            'includedData': 'attributes,identifiers,images,productTypes,salesRanks,summaries,variations'
        }
        response = await self.request(method, url, params=query)

        logger.info('action=get_catalog_item status=done')
        return response

    async def search_catalog_items_v2022_04_01(self, identifiers: List[str], id_type: str) -> dict:
        logger.info('action=search_catalog_items_v2022_04_01 status=run')

        method = "GET"
        path = "/catalog/2022-04-01/items"
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'identifiers': ','.join(identifiers),
            'identifiersType': id_type,
            'marketplaceIds': self.marketplace_id,
            'includedData': 'attributes,dimensions,identifiers,productTypes,relationships,salesRanks,summaries',
            'pageSize': 20,
        }
        response = await self.request(method, url, query)

        logger.info('action=search_catalog_items_v2022_04_01 status=done')
        return response

    async def get_item_offers_batch(self, asin_list: List, item_condition: str='NEW', customer_type: str='Consumer') -> dict:
        logger.info('action=get_item_offers_batch status=run')

        if len(asin_list) > 20:
            raise TooMatchParameterException

        request_list = []
        for asin in asin_list:
            request_list.append({
                'uri': f'/products/pricing/v0/items/{asin}/offers',
                'method': 'GET',
                'MarketplaceId': self.marketplace_id,
                'ItemCondition': item_condition,
                'CustomerType': customer_type,
            })

        method = 'POST'
        path = '/batches/products/pricing/v0/itemOffers'
        url = urllib.parse.urljoin(ENDPOINT, path)
        body = {
            'requests' : request_list,
        }
        response = await self.request(method, url, body=body)

        logger.info('action=get_item_offers_batch status=done')
        return response


class SPAPIJsonParser(object):

    @staticmethod
    def parse_get_competitive_pricing(response: dict) -> dict:
        logger.info('action=parse_get_competitive_pricing status=run')

        products = []

        for payload in response.get('payload'):
            asin = payload['ASIN']
            try:
                price = round(float(payload['Product']['CompetitivePricing']['CompetitivePrices'][0]['Price']['LandedPrice']['Amount']))
            except (IndexError, KeyError) as ex:
                logger.error(f"{asin} hasn't landedprice error={ex}")
                price = -1

            try:
                ranking = round(float(payload['Product']['SalesRankings'][0]['Rank']))
                category_id = payload['Product']['SalesRankings'][0]['ProductCategoryId']
                if re.fullmatch('[\d]+', category_id):
                    raise NotRankingException
            except (NotRankingException, IndexError, KeyError) as ex:
                logger.error(f"{asin} hasn't ranking error={ex}")
                ranking = -1

            products.append({'asin': asin, 'price': price, 'ranking': ranking})

        logger.info('action=parse_get_competitive_pricing status=done')
        return products

    @staticmethod
    def parse_get_item_offers(response: dict) -> dict|None:
        logger.info('action=parse_get_item_offers status=run')

        try:
            asin = response['payload']['ASIN']
        except KeyError as ex:
            logger.error(ex)
            logger.error(response)
            return None

        try:
            price = int(response['payload']['Summary']['LowestPrices'][0]['LandedPrice']['Amount'])
            ranking = response['payload']['Summary']['SalesRankings'][0]['Rank']
        except (IndexError, KeyError) as ex:
            logger.error(f"{asin} hasn't data")
            return None

        logger.info('action=parse_get_item_offers status=done')
        return {'asin': asin, 'price': price, 'ranking': ranking}

    @staticmethod
    def parse_list_catalog_items(response: dict) -> list[dict]:
        logger.info('action=parse_list_catalog_items status=run')

        products = []
        try:
            items = response['payload']['Items']
        except KeyError as ex:
            logger.error(ex)
            return products

        for item in items:
            asin = item['Identifiers']['MarketplaceASIN']['ASIN']
            title = item['AttributeSets'][0]['Title']

            try:
                quantity = item['AttributeSets'][0]['PackageQuantity']
            except KeyError as ex:
                logger.error(ex)
                quantity = 1
            products.append({'asin': asin, 'quantity': quantity, 'title': title})
            
        logger.info('action=parse_list_catalog_items status=done')
        return products

    @staticmethod
    def parse_get_my_fees_estimate_for_asin(response: dict, amount: int=10000, default_fee_rate: float=0.1, default_ship_fee: int=500) -> dict:
        logger.info('action=parse_get_my_fees_estimate_for_asin status=run')

        fee_type_dict = {}

        try:
            asin = response['payload']['FeesEstimateResult']['FeesEstimateIdentifier']['SellerInputIdentifier']
        except KeyError as ex:
            logger.error(response)
            time.sleep(1)
            raise QuotaException

        try:
            fees = response['payload']['FeesEstimateResult']['FeesEstimate']['FeeDetailList']
        except KeyError as ex:
            logger.info(ex)
            return {'asin': asin, 'fee_rate': default_fee_rate, 'ship_fee': default_ship_fee}
        
        for fee in fees:
            try:
                fee_type_dict[fee['FeeType']] = fee['FeeAmount']['Amount']
            except KeyError as ex:
                logger.info(ex)

        fee_rate = round(fee_type_dict.get('ReferralFee') / amount, 2)
        ship_fee = fee_type_dict.get('FBAFees', default_ship_fee)

        logger.info('action=parse_get_my_fees_estimate_for_asin status=done')
        return {'asin': asin, 'fee_rate': fee_rate, 'ship_fee': ship_fee}

    @staticmethod
    @logger_decorator
    def parse_search_catalog_items_v2022_04_01(response: dict) -> List[dict]:
        products = []
        
        items = response['items']
        for item in items:
            asin = item['asin']

            try:
                item_name_list = item['attributes']['item_name']
                for item_name in item_name_list:
                    title = item_name['value']
                    if title:
                        break
            except KeyError as ex:
                logger.error({'message': 'item name is None', 'error': ex})
                continue

            try:
                unit_count_list = item['attributes']['unit_count']
                for unit_count in unit_count_list:
                    quantity = unit_count['value']
                    if quantity:
                        break
            except KeyError as ex:
                logger.info({'message': 'unit info is None', 'error': ex})
                quantity = 1

            products.append({'asin': asin, 'quantity': int(float(quantity)), 'title': title})
        return products

    @classmethod
    @logger_decorator
    def parse_get_item_offers_batch(cls, response: dict) -> List[dict]:

        products = []
        responses = response['responses']
        for response in responses:
            try:
                result = cls.parse_get_item_offers(response['body'])
                products.append(result)
            except KeyError as ex:
                logger.error({'message':ex})
                continue

        return products


class NotRankingException(Exception):
    pass


class QuotaException(Exception):
    pass

class TooMatchParameterException(Exception):
    pass
