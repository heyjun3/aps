import datetime
import hmac
import hashlib
import json
import os
import time
import urllib.parse
import re

import redis
import requests

import settings
import log_settings


logger = log_settings.get_logger(__name__)
redis_client = redis.Redis(
        host=settings.REDIS_HOST,
        port=settings.REDIS_PORT,
        db=settings.REDIS_DB,
    )
ENDPOINT = 'https://sellingpartnerapi-fe.amazon.com'


def request(req: requests.Request) -> requests.Response:
    for _ in range(60):
        try:
            session = requests.Session()
            response = session.send(req.prepare())
            if response.status_code == 200 or response is not None:
                return response
            else:
                raise Exception
        except Exception as ex:
            logger.error(f'action=request error={ex}')
            time.sleep(60)


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

    def get_spapi_access_token(self, timeout_sec: int = 3500):

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
        req = requests.Request(method='POST', url=URL, params=params, headers=headers)
        response = request(req)

        access_token = response.json().get('access_token')
        
        if access_token is None:
            logger.error(response.text)
            raise Exception

        redis_client.set('access_token', access_token, ex=timeout_sec)

        return access_token

    def create_authorization_headers(self, req: requests.Request) -> requests.Request:
        region = 'us-west-2'
        service = 'execute-api'
        algorithm = 'AWS4-HMAC-SHA256'
        signed_headers = 'host;user-agent;x-amz-access-token;x-amz-date'
        user_agent = 'My SPAPI Client tool /1.0(Language=python/3.10)'

        host = urllib.parse.urlparse(ENDPOINT).netloc
        canonical_uri = urllib.parse.urlparse(req.url).path
        body = ''

        if  req.json:
            body = json.dumps(req.json)

        utcnow = datetime.datetime.utcnow()
        amz_date = utcnow.strftime('%Y%m%dT%H%M%SZ')
        datestamp = utcnow.strftime('%Y%m%d')
        amz_access_token = self.get_spapi_access_token()
        canonical_header_values = [host, user_agent, amz_access_token, f'{amz_date}\n']
        
        canonical_headers = '\n'.join([f'{head}:{value}' for head, value in zip(signed_headers.split(';'), canonical_header_values)])
        payload_hash = hashlib.sha256(body.encode('utf-8')).hexdigest()
        canonical_querystring = ''

        if req.params:
            canonical_querystring = urllib.parse.urlencode(sorted(req.params.items(), key=lambda x: (x[0], x[1])))

        canonical_request = '\n'.join([req.method, canonical_uri, canonical_querystring, canonical_headers, signed_headers, payload_hash])
        credential_scope = os.path.join(datestamp, region, service, 'aws4_request')
        signing_key = self.get_signature_key(self.aws_secret_key, datestamp, region, service)
        string_to_sign = '\n'.join([algorithm, amz_date, credential_scope, hashlib.sha256(canonical_request.encode('utf-8')).hexdigest()])
        
        signature = hmac.new(signing_key, (string_to_sign).encode('utf-8'), hashlib.sha256).hexdigest()
        authorization_header = f'{algorithm} Credential={self.aws_access_key}/{credential_scope}, SignedHeaders={signed_headers}, Signature={signature}'
        
        req.headers = {
            'host': urllib.parse.urlparse(ENDPOINT).netloc,
            'user-agent': user_agent,
            'x-amz-date': amz_date,
            'Authorization': authorization_header,
            'x-amz-access-token': amz_access_token,
        }

        return req

    def get_my_fees_estimate_for_asin(self, asin: str, price: int, is_fba: bool = True, currency_code: str = 'JPY'):
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
        req = requests.Request(method=method, url=url, json=body)
        req = self.create_authorization_headers(req)
        response = request(req)

        return response

    def get_pricing(self, asin_list: list, item_type: str='Asin') -> requests.Response:
        method = 'GET'
        path = '/products/pricing/v0/price'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'Asins': ','.join(asin_list),
            'ItemType': item_type,
            'MarketplaceId': self.marketplace_id,
        }
        req = requests.Request(method=method, url=url, params=query)
        req = self.create_authorization_headers(req)
        response = request(req)

        return response

    def get_competitive_pricing(self, asin_list: list, item_type: str='Asin') -> requests.Response:
        logger.info('action=get_competitive_pricing status=run')

        method = 'GET'
        path = '/products/pricing/v0/competitivePrice'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'MarketplaceId': self.marketplace_id,
            'Asins': ','.join(asin_list),
            'ItemType': item_type,
        }
        req = requests.Request(method=method, url=url, params=query)
        req = self.create_authorization_headers(req)
        response = request(req)

        logger.info('action=get_competitive_pricing status=done')
        return response

    def get_item_offers(self, asin: str, item_condition: str='New') -> requests.Response:
        logger.info('action=get_item_offers status=run')

        method = 'GET'
        path = f'/products/pricing/v0/items/{asin}/offers'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'MarketplaceId': self.marketplace_id,
            'ItemCondition': item_condition,
        }
        req = requests.Request(method=method, url=url, params=query)
        req = self.create_authorization_headers(req)
        response = request(req)

        logger.info('action=get_item_offers status=done')
        return response

    def search_catalog_items(self, jan_list: list) -> requests.Response:
        logger.info('action=search_catalog_items status=run')

        method = 'GET'
        path = '/catalog/2020-12-01/items'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'keywords': ','.join(jan_list),
            'marketplaceIds': self.marketplace_id,
            'includedData': 'identifiers,images,productTypes,salesRanks,summaries,variations'
        }
        req = requests.Request(method=method, url=url, params=query)
        req = self.create_authorization_headers(req)
        response = request(req)

        logger.info('action=search_catalog_items status=done')
        return response

    def list_catalog_items(self, jan: str) -> requests.Request:
        logger.info('action=list_catalog_items status=run')

        method = 'GET'
        path = '/catalog/v0/items'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'MarketplaceId': self.marketplace_id,
            'JAN': jan,
        }
        req = requests.Request(method=method, url=url, params=query)
        req = self.create_authorization_headers(req)
        response = request(req)

        logger.info('action=list_catalog_items status=done')
        return response

    def get_catalog_item(self, asin: str) -> requests.Response:
        logger.info('action=get_catalog_item status=run')

        method = 'GET'
        path = f'/catalog/2020-12-01/items/{asin}'
        url = urllib.parse.urljoin(ENDPOINT, path)
        query = {
            'marketplaceIds': self.marketplace_id,
            'includedData': 'attributes,identifiers,images,productTypes,salesRanks,summaries,variations'
        }
        req = requests.Request(method=method, url=url, params=query)
        req = self.create_authorization_headers(req)
        response = request(req)

        logger.info('action=get_catalog_item status=done')
        return response


class SPAPIJsonParser(object):

    @staticmethod
    def parse_get_competitive_pricing(response: json) -> dict:
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
    def parse_get_item_offers(response: json) -> dict|None:
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
    def parse_list_catalog_items(response: json) -> list[dict]:
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
            try:
                price = item['AttributeSets'][0]['ListPrice']['Amount']
            except KeyError as ex:
                logger.error(ex)
                price = None
            products.append({'asin': asin, 'quantity': quantity, 'title': title, 'price': price})
            
        logger.info('action=parse_list_catalog_items status=done')
        return products


class NotRankingException(Exception):
    pass