import datetime
import hmac
import hashlib
import json
import time
import urllib.parse


import requests

import settings
import log_settings


logger = log_settings.get_logger(__name__)

    
ENDPOINT = 'https://sellingpartnerapi-fe.amazon.com'


def request(method: str, url: str, params: dict = None, json: dict = None, headers: dict = None):
    for _ in range(69):
        try:
            response = requests.request(method, url=url, params=params, json=json, headers=headers)
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

    def get_spapi_access_token(self):

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

        response = request('POST', URL, params=params, headers=headers)
        access_token = response.json().get('access_token')
        
        if access_token is None:
            logger.error(response.text)
            return None

        return access_token

    def create_signature(self, method, canonical_uri, access_token: str, body, query: str = ''):
        region = 'us-west-2'
        service = 'execute-api'
        user_agent = 'My SPAPI Client tool /1.0(Language=python/3.10'
        algorithm = 'AWS4-HMAC-SHA256'
        signed_headers = 'host;user-agent;x-amz-access-token;x-amz-date'

        host = urllib.parse.urlparse(ENDPOINT).netloc

        t = datetime.datetime.utcnow()
        amz_date = t.strftime('%Y%m%dT%H%M%SZ') 
        datestamp = t.strftime('%Y%m%d')

        canonical_headers = 'host:' + host + '\n' + 'user-agent:' + user_agent + '\n' + 'x-amz-access-token:' + access_token + '\n' + 'x-amz-date:' + amz_date + '\n'
        payload_hash = hashlib.sha256(json.dumps(body)).encode('utf-8').hexdigest()

        canonical_querystring = urllib.parse.urlencode(query)

        canonical_request = method + '\n' + canonical_uri + '\n' + canonical_querystring + '\n' + canonical_headers + '\n' + signed_headers + '\n' + payload_hash
        credential_scope = datestamp + '/' + region + '/' + service + '/' + 'aws4_request'
        signing_key = self.get_signature_key(self.aws_secret_key, datestamp, region, service)
        string_to_sign = algorithm + '\n' +  amz_date + '\n' +  credential_scope + '\n' +  hashlib.sha256(canonical_request.encode('utf-8')).hexdigest()
        
        signature = hmac.new(signing_key, (string_to_sign).encode('utf-8'), hashlib.sha256).hexdigest()
        authorization_header = algorithm + ' ' + 'Credential=' + self.aws_access_key + '/' + credential_scope + ', ' +  'SignedHeaders=' + signed_headers + ', ' + 'Signature=' + signature
        headers = {
            'host': host,
            'user-agent': 'My SPAPI Client tool /1.0(Language=python/3.10',
            'x-amz-access-token': access_token,
            'x-amz-date': amz_date,
            'Authorization': authorization_header,
        }
        return headers

    def get_my_fees_estimate_for_asin(self, asin: str, price: int, is_fba: bool = True, currency_code: str = 'JPY'):
        logger.info('action=get_my_fees_estimate_for_asin status=run')

        method = 'POST'
        path = f'/products/fees/v0/items/{asin}/feesEstimate'
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
        # query = {
        #     'Asins': asin,
        #     'ItemType': 'Asin',
        #     'MarketplaceId': settings.MARKETPLACEID,

        # }
        access_token = self.get_spapi_access_token()
        url = urllib.parse.urljoin(ENDPOINT, path)
        headers = self.create_signature(method, url, access_token=access_token, body=body)
        response = requests.post(url, json=body, headers=headers)
        print(response.text)
