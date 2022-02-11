import urllib
import datetime

import requests

import settings
import log_settings


logger = log_settings.get_logger(__name__)


def get_my_fees_estimate_for_asin():
    ENDPOINT = 'https://sellingpartnerapi-fe.amazon.com'
    asin = 'B074X15BG5'
    price = 28500
    currency_code = 'JPY'
    is_fba = True
    marketplace_id = settings.MARKETPLACEID

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
            'MarketplaceId': marketplace_id, 
            }
        }
    # query = urllib.parse.urlencode(body)
    # headers = {
    #     'host': urllib.parse.urlparse(ENDPOINT).netloc,
    #     'user-agent': 'My SPAPI Client tool /1.0(Language=python/3.10',
    #     'x-amz-access-token': get_spapi_access_token(),
    #     'x-amz-date': datetime.datetime.now().strftime('%Y%m%dT%H%M%SZ')
    # }
    # print(query)
    query = urllib.parse.urlencode(body)
    print(query)
    url = urllib.parse.urlparse(query)
    print(url)
    url = urllib.parse.urljoin(ENDPOINT, f'/products/fees/v0/items/{asin}/feesEstimate')
    # response = requests.get(url, params=body)
    # print(response.url)


import hmac
import hashlib

def sign(key, msg):
    return hmac.new(key, msg.encode('utf-8'), hashlib.sha256).digest()


def get_signature_key(key, dateStamp, regionName, serviceName):
    kDate = sign(('AWS4' + key).encode('utf-8'), dateStamp)
    kRegion = sign(kDate, regionName)
    kService = sign(kRegion, serviceName)
    kSigning = sign(kService, 'aws4_request')
    return kSigning


def create_signature():
    asin = 'B074X15BG5'
    region = 'us-west-2'
    service = 'execute-api'
    method = 'POST'
    canonical_uri = f'/products/fees/v0/items/{asin}/feesEstimate'
    secret_key = settings.AWS_SECRET_KEY

    t = datetime.datetime.utcnow()
    amz_date = t.strftime('%Y%m%dT%H%M%SZ') 
    datestamp = t.strftime('%Y%m%d')
    algorithm = 'AWS4-HMAC-SHA256'
    canonical_request = method + '\n' + canonical_uri + '\n' + canonical_querystring + '\n' + canonical_headers + '\n' + signed_headers + '\n' + payload_hash
    credential_scope = datestamp + '/' + region + '/' + service + '/' + 'aws4_request'
    signing_key = get_signature_key(secret_key, datestamp, region, service)
    string_to_sign = algorithm + '\n' +  amz_date + '\n' +  credential_scope + '\n' +  hashlib.sha256(canonical_request.encode('utf-8')).hexdigest()
    
    signature = hmac.new(signing_key, (string_to_sign).encode('utf-8'), hashlib.sha256).hexdigest()
    pass


def get_spapi_access_token():

    URL = 'https://api.amazon.com/auth/o2/token'
    headers = {
    'Host': 'api.amazon.com',
    'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
    }
    params = {
    'grant_type': 'refresh_token',
    'refresh_token': settings.REFRESH_TOKEN,
    'client_id': settings.CLIENT_ID,
    'client_secret': settings.ClIENT_TOKEN,
    }

    response = requests.post(URL, params=params, headers=headers)
    access_token = response.json().get('access_token')
    
    if access_token is None:
        logger.error(response.text)
        return None

    return access_token