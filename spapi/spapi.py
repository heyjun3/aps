import datetime
import hmac
import hashlib
import json
import os
import time
import urllib.parse
import re
import asyncio
from functools import partial
from typing import List
from typing import Callable

import redis
import aiohttp

import settings
import log_settings


logger = log_settings.get_logger(__name__)


def logger_decorator(func: Callable) -> Callable:
    def _logger_decorator(*args, **kwargs):
        logger.info({'action': func.__name__, 'status': 'run'})
        result = func(*args, **kwargs)
        logger.info({'action': func.__name__, 'status': 'done'})
        return result
    return _logger_decorator


async def request(method: str, url: str, params: dict=None, headers: dict=None, body: dict=None) -> aiohttp.ClientResponse:
    for _ in range(60):
        async with aiohttp.request(method, url, params=params, headers=headers, json=body) as response:
            response_json = await response.json()
            if response.status == 200 and response is not None or response.status == 400:
                response_json = await response.json()
                return response_json
            elif response.status == 429:
                logger.error(response_json)
                await asyncio.sleep(2)
            else:
                logger.error(response_json)
                await asyncio.sleep(10)
                

class SPAPI(object):

    def __init__(self):
        self.refresh_toke = settings.REFRESH_TOKEN
        self.client_id = settings.CLIENT_ID
        self.client_secret = settings.CLIENT_SECRET
        self.aws_secret_key = settings.AWS_SECRET_KEY
        self.aws_access_key = settings.AWS_ACCESS_ID
        self.marketplace_id = settings.MARKETPLACEID
        self.seller_id = settings.SELLER_ID
        self.redis_client = redis.Redis(
            host=settings.REDIS_HOST,
            port=settings.REDIS_PORT,
            db=settings.REDIS_DB,
            password=settings.REDIS_PASSWORD,
        )

    def sign(self, key, msg):
        return hmac.new(key, msg.encode('utf-8'), hashlib.sha256).digest()

    def get_signature_key(self, key, dateStamp, regionName, serviceName):
        kDate = self.sign(('AWS4' + key).encode('utf-8'), dateStamp)
        kRegion = self.sign(kDate, regionName)
        kService = self.sign(kRegion, serviceName)
        kSigning = self.sign(kService, 'aws4_request')
        return kSigning

    async def get_spapi_access_token(self, timeout_sec: int = 3500) -> str:

        access_token = self.redis_client.get('access_token')
        if access_token is not None:
            return access_token.decode()

        method, url, params, headers = self._get_ephemeral_access_token()
        response = await request(method, url, params, headers)

        access_token = response.get('access_token')
        if access_token is None:
            text = await response.text()
            logger.error(text)
            raise Exception

        self.redis_client.set('access_token', access_token, ex=timeout_sec)

        return access_token

    def _get_ephemeral_access_token(self) -> tuple[str, str, dict, dict]:

        method = 'POST'
        url = 'https://api.amazon.com/auth/o2/token'

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
        return (method, url, params, headers)

    def create_authorization_headers(self, amz_access_token: str, method: str, url: str, params: dict={}, body: dict={}) -> dict:
        region = 'us-west-2'
        service = 'execute-api'
        algorithm = 'AWS4-HMAC-SHA256'
        signed_headers = 'host;user-agent;x-amz-access-token;x-amz-date'
        user_agent = 'My SPAPI Client tool /1.0(Language=python/3.10)'

        host = urllib.parse.urlparse(settings.ENDPOINT).netloc
        canonical_uri = urllib.parse.urlparse(url).path

        body_json = ''
        if body:
            body_json = json.dumps(body)

        utcnow = datetime.datetime.utcnow()
        amz_date = utcnow.strftime('%Y%m%dT%H%M%SZ')
        datestamp = utcnow.strftime('%Y%m%d')
        canonical_header_values = [host, user_agent, amz_access_token, f'{amz_date}\n']
        
        canonical_headers = '\n'.join([f'{head}:{value}' for head, value in zip(signed_headers.split(';'), canonical_header_values)])
        payload_hash = hashlib.sha256(body_json.encode('utf-8')).hexdigest()

        canonical_querystring = ''
        if params:
            canonical_querystring = urllib.parse.urlencode(sorted(params.items(), key=lambda x: (x[0], x[1])))

        canonical_request = '\n'.join([method, canonical_uri, canonical_querystring, canonical_headers, signed_headers, payload_hash])
        credential_scope = os.path.join(datestamp, region, service, 'aws4_request')
        signing_key = self.get_signature_key(self.aws_secret_key, datestamp, region, service)
        string_to_sign = '\n'.join([algorithm, amz_date, credential_scope, hashlib.sha256(canonical_request.encode('utf-8')).hexdigest()])
        
        signature = hmac.new(signing_key, (string_to_sign).encode('utf-8'), hashlib.sha256).hexdigest()
        authorization_header = f'{algorithm} Credential={self.aws_access_key}/{credential_scope}, SignedHeaders={signed_headers}, Signature={signature}'
        
        headers = {
            'host': urllib.parse.urlparse(settings.ENDPOINT).netloc,
            'user-agent': user_agent,
            'x-amz-date': amz_date,
            'Authorization': authorization_header,
            'x-amz-access-token': amz_access_token,
        }

        return headers
    
    async def _request(self, func: Callable) -> dict:
        method, url, params, body = func()
        access_token = await self.get_spapi_access_token()
        headers = self.create_authorization_headers(access_token, method, url, params, body)
        response = await request(method, url, params=params, body=body, headers=headers)
        return response

    async def get_my_fees_estimate_for_asin(self, asin: str, price: int=10000, is_fba: bool=True) -> dict:
        return await self._request(partial(self._get_my_fees_estimate_for_asin, asin, price, is_fba))

    async def get_pricing(self, asin_list: List[str], item_type: str='Asin') -> dict:
        return await self._request(partial(self._get_pricing, asin_list, item_type))

    async def get_competitive_pricing(self, asin_list: List[str], item_type: str='Asin') -> dict:
        return await self._request(partial(self._get_competitive_pricing, asin_list, item_type))

    async def get_item_offers(self, asin: str, item_condition: str='New') -> dict:
        return await self._request(partial(self._get_item_offers, asin, item_condition))

    async def search_catalog_items(self, jan_list: List[str]) -> dict:
        return await self._request(partial(self._search_catalog_items, jan_list))

    async def list_catalog_items(self, jan: str) -> dict:
        return await self._request(partial(self._list_catalog_items, jan))

    async def get_catalog_item(self, asin: str) -> dict:
        return await self._request(partial(self._get_catalog_item, asin))

    async def search_catalog_items_v2022_04_01(self, identifiers: List[str], id_type: str) -> dict:
        return await self._request(partial(self._search_catalog_items_v2022_04_01, identifiers, id_type))

    async def get_item_offers_batch(self, asin_list: List[str], item_condition: str='NEW', customer_type: str='Consumer') -> dict:
        return await self._request(partial(self._get_item_offers_batch, asin_list, item_condition, customer_type))

    async def get_my_fees_estimates(self, asin_list: List[str], id_type: str='ASIN', price_amount: int=10000) -> dict:
        return await self._request(partial(self._get_my_fees_estimates, asin_list, id_type, price_amount))

    def _get_my_fees_estimate_for_asin(self, asin: str, price: int=10000, is_fba: bool = True, currency_code: str = 'JPY') -> tuple[str, str, None, dict]:
        logger.info('action=get_my_fees_estimate_for_asin status=run')

        method = 'POST'
        path = f'/products/fees/v0/items/{asin}/feesEstimate'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = None
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
        return (method, url, query, body)

    def _get_pricing(self, asin_list: list, item_type: str='Asin') -> tuple[str, str, dict, None]:
        method = 'GET'
        path = '/products/pricing/v0/price'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            'Asins': ','.join(asin_list),
            'ItemType': item_type,
            'MarketplaceId': self.marketplace_id,
        }
        body = None

        return (method, url, query, body)

    def _get_competitive_pricing(self, asin_list: list, item_type: str='Asin') -> tuple[str, str, dict, None]:
        logger.info('action=get_competitive_pricing status=run')

        method = 'GET'
        path = '/products/pricing/v0/competitivePrice'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            'MarketplaceId': self.marketplace_id,
            'Asins': ','.join(asin_list),
            'ItemType': item_type,
        }
        body = None

        logger.info('action=get_competitive_pricing status=done')
        return (method, url, query, body)

    def _get_item_offers(self, asin: str, item_condition: str='New') -> tuple[str, str, dict, None]:
        logger.info('action=get_item_offers status=run')

        method = 'GET'
        path = f'/products/pricing/v0/items/{asin}/offers'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            'MarketplaceId': self.marketplace_id,
            'ItemCondition': item_condition,
        }
        body = None

        logger.info('action=get_item_offers status=done')
        return (method, url, query, body)

    def _search_catalog_items(self, jan_list: list) -> tuple[str, str, dict, None]:
        logger.info('action=search_catalog_items status=run')

        method = 'GET'
        path = '/catalog/2020-12-01/items'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            'keywords': ','.join(jan_list),
            'marketplaceIds': self.marketplace_id,
            'includedData': 'identifiers,images,productTypes,salesRanks,summaries,variations'
        }
        body = None

        logger.info('action=search_catalog_items status=done')
        return (method, url, query, body)

    def _list_catalog_items(self, jan: str) -> tuple[str, str, dict, None]:
        logger.info('action=list_catalog_items status=run')

        method = 'GET'
        path = '/catalog/v0/items'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            'MarketplaceId': self.marketplace_id,
            'JAN': jan,
        }
        body = None 

        logger.info('action=list_catalog_items status=done')
        return (method, url, query, body)

    def _get_catalog_item(self, asin: str) -> tuple[str, str, dict, None]:
        logger.info('action=get_catalog_item status=run')

        method = 'GET'
        path = f'/catalog/2020-12-01/items/{asin}'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            'marketplaceIds': self.marketplace_id,
            'includedData': 'attributes,identifiers,images,productTypes,salesRanks,summaries,variations'
        }
        body = None

        logger.info('action=get_catalog_item status=done')
        return (method, url, query, body)

    def _search_catalog_items_v2022_04_01(self, identifiers: List[str], id_type: str) -> tuple[str, str, dict, None]:
        logger.info('action=search_catalog_items_v2022_04_01 status=run')

        if (len(identifiers) > 20):
            raise TooMatchParameterException

        method = "GET"
        path = "/catalog/2022-04-01/items"
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = {
            'identifiers': ','.join(identifiers),
            'identifiersType': id_type,
            'marketplaceIds': self.marketplace_id,
            'includedData': 'attributes,dimensions,identifiers,productTypes,relationships,salesRanks,summaries',
            'pageSize': 20,
        }
        body = None

        logger.info('action=search_catalog_items_v2022_04_01 status=done')
        return (method, url, query, body)

    def _get_item_offers_batch(self, asin_list: List, item_condition: str='NEW', customer_type: str='Consumer') -> tuple[str, str, None, dict]:
        logger.info('action=get_item_offers_batch status=run')

        if len(asin_list) > 20:
            raise TooMatchParameterException

        request_list = []
        for asin in list(set(asin_list)):
            request_list.append({
                'uri': f'/products/pricing/v0/items/{asin}/offers',
                'method': 'GET',
                'MarketplaceId': self.marketplace_id,
                'ItemCondition': item_condition,
                'CustomerType': customer_type,
            })

        method = 'POST'
        path = '/batches/products/pricing/v0/itemOffers'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = None
        body = {
            'requests' : request_list,
        }

        logger.info('action=get_item_offers_batch status=done')
        return (method, url, query, body)

    def _get_my_fees_estimates(self, asin_list: List, id_type: str='ASIN', price_amount: int=10000) -> tuple[str, str, None, dict]:
        logger.info('action=get_my_fees_estimates status=run')

        if len(asin_list) > 20:
            raise TooMatchParameterException

        method = 'POST'
        path = '/products/fees/v0/feesEstimate'
        url = urllib.parse.urljoin(settings.ENDPOINT, path)
        query = None

        body = []
        for asin in asin_list:
            body.append({
                'FeesEstimateRequest': {
                    'MarketplaceId': self.marketplace_id,
                    'IsAmazonFulfilled': True,
                    'PriceToEstimateFees': {
                        'ListingPrice': {
                            'CurrencyCode': 'JPY',
                            'Amount': price_amount,
                        },
                    },
                    'Identifier': asin,
                    'OptionalFulfillmentProgram': 'FBA_CORE',
                },
                'IdType': id_type,
                'IdValue': asin,
            })
        
        logger.info('action=get_my_fees_estimates status=done')
        return (method, url, query, body)


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
                if re.fullmatch('[0-9]+', category_id):
                    raise NotRankingException
            except (NotRankingException, IndexError, KeyError) as ex:
                logger.error(f"{asin} hasn't ranking error={ex}")
                ranking = -1

            products.append({'asin': asin, 'price': price, 'ranking': ranking})

        logger.info('action=parse_get_competitive_pricing status=done')
        return products

    @staticmethod
    def parse_get_item_offers(response: dict) -> dict:
        logger.info('action=parse_get_item_offers status=run')
            
        asin = response['payload']['ASIN']
        try:
            new_buy_box = list(filter(lambda x: x['condition'] == 'new', response['payload']['Summary']['BuyBoxPrices']))
            price = int(new_buy_box[0]['LandedPrice']['Amount'])
        except (IndexError, KeyError) as ex:
            logger.info('new buy box information is None')
            try:
                lowest_offer = response['payload']['Offers'][0]
                price = int(lowest_offer['ListingPrice']['Amount'])
                shipping_cost = int(lowest_offer['Shipping']['Amount'])
                price += shipping_cost
            except (IndexError, KeyError) as ex:
                logger.info('lowest price offer is None')
                price = -1

        try:
            ranking = response['payload']['Summary']['SalesRankings'][0]['Rank']
        except (IndexError, KeyError) as ex:
            logger.error(f"error={ex}")
            ranking = -1

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
                identifiers = item['identifiers']
                for identifier in identifiers:
                    jan = identifier['identifiers'][0]['identifier']
                    if jan:
                        break
            except (KeyError, IndexError) as ex:
                logger.error({'message': 'jan is None', 'error': ex})
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

            products.append({'asin': asin, 'quantity': int(float(quantity)), 'title': title, 'jan': jan})
        return products

    @classmethod
    @logger_decorator
    def parse_get_item_offers_batch(cls, response: dict) -> List[dict]:

        products = []
        responses = response['responses']
        for response in responses:
            try:
                status_code = response.get('status').get('statusCode')
                if status_code == 200:
                    result = cls.parse_get_item_offers(response['body'])
                    products.append(result)
                else:
                    asin = response.get('request').get('Asin')
                    products.append({'asin': asin, 'price': -1, 'ranking': -1})
            except KeyError as ex:
                logger.error({'message':ex})
                continue

        return products

    @staticmethod
    @logger_decorator
    def parse_get_my_fees_estimates(response: dict, default_fee_rate: float=0.1, default_ship_fee: int=500) -> List[dict]:
        products = []

        for product in response:
            asin = product["FeesEstimateIdentifier"]['IdValue']
            amount = product['FeesEstimateIdentifier']['PriceToEstimateFees']["ListingPrice"]["Amount"]

            try:
                fee_detail_list = product['FeesEstimate']["FeeDetailList"]
                fee = [fee_detail for fee_detail in fee_detail_list if fee_detail.get('FeeType') == "ReferralFee"]
                if not fee:
                    fee_rate = default_fee_rate
                else:
                    fee = fee[0]['FeeAmount']["Amount"]
                    fee_rate = round(int(fee) / int(amount), 2)
            except (KeyError, IndexError) as ex:
                logger.error({'message': ex})
                fee_rate = default_fee_rate

            try:
                fee_detail_list = product['FeesEstimate']["FeeDetailList"]
                ship_fee = [fee_detail for fee_detail in fee_detail_list if fee_detail.get('FeeType') == "FBAFees"]
                if not ship_fee:
                    ship_fee = default_ship_fee
                else:
                    ship_fee = ship_fee[0]["FeeAmount"]["Amount"]
            except (KeyError, IndexError) as ex:
                logger.error({'message': ex})
                ship_fee = default_ship_fee

            products.append({'asin': asin, 'fee_rate': fee_rate, 'ship_fee': ship_fee})
        
        return products


class NotRankingException(Exception):
    pass


class QuotaException(Exception):
    pass

class TooMatchParameterException(Exception):
    pass
