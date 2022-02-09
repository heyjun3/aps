from datetime import datetime
import urllib.parse
import hmac
import six
import hashlib
import base64
import time
from multiprocessing import Process
import logging.config
import os
import pathlib
import shutil
import urllib.parse

import requests
import xml.etree.ElementTree as et
import pandas as pd
from lxml import etree
from bs4 import BeautifulSoup

from mws.models import MWS
import settings


logger = logging.getLogger(__name__)


def datetime_encode(dt):
    return dt.strftime('%Y-%m-%dT%H:%M:%SZ')


def request_api(url, data=None):
    logger.info('action=request_api status=run')

    for _ in range(60):
        try:
            response = requests.post(url, data=data, timeout=30.0)
            if not response.status_code == 200 or response is None:
                logger.error(response.text)
                raise Exception
            logger.info('action=request_api status=done')
            return response
        except Exception as e:
            logger.error(f'action=request_api error={e}')
            time.sleep(60)


class AmazonClient:
    def __init__(self):
        self.secret_key = settings.SECRET_KEY
        self.domain = settings.DOMAIN
        self.endpoint = settings.ENDPOINT
        self.data = {
            'AWSAccessKeyId': settings.ACCESS_ID,
            'SellerId': settings.SELLER_ID,
            'SignatureMethod': 'HmacSHA256',
            'SignatureVersion': '2',
            'Version': '2011-10-01',
        }
        self.tag_xmlns = settings.XMLNS
        self.tag_ns2 = settings.NS2

    def request(self, data_dict):
        logger.info('action=request status=run')

        url = self.create_request_url(data_dict)
        response = request_api(url)
        response_text = response.content.decode()
        response_xml = et.fromstring(response_text)
        return response_xml

    def create_request_url(self, data_dict):
        logger.debug('action=create_request_url status=run')

        timestamp = datetime_encode(datetime.utcnow())
        data_dict['Timestamp'] = timestamp

        query_string = '&'.join('{}={}'.format(
            n, urllib.parse.quote(v, safe='')) for n, v in sorted(data_dict.items()))
        canonical = "{}\n{}\n{}\n{}".format(
            'POST', self.domain, self.endpoint, query_string)

        h = hmac.new(
            six.b(self.secret_key),
            six.b(canonical), hashlib.sha256)

        signature = urllib.parse.quote(base64.b64encode(h.digest()), safe='')

        url = 'https://{}{}?{}&Signature={}'.format(
            self.domain, self.endpoint, query_string, signature)

        logger.debug('action=create_request_url status=done')
        return url

    def get_matching_product_for_id(self, products_dict: dict, interval_sec: int = 1):
        logger.info('action=get_matching_product_for_id status=run')
        mws_object_list = []
        
        products = list(products_dict.keys())
        data_dict = dict(self.data)
        data_dict['MarketplaceId'] = settings.MARKETPLACEID
        data_dict['IdType'] = 'JAN'
        data_dict['Action'] = 'GetMatchingProductForId'
        for index, jan in enumerate(products):
            data_dict[f'IdList.Id.{str(index+1)}'] = str(jan)

        url = self.create_request_url(data_dict=data_dict)
        response = request_api(url)
        response.encoding = response.apparent_encoding
        time.sleep(interval_sec)
        tree = etree.fromstring(response.text)

        for result in tree.findall(".//GetMatchingProductForIdResult", tree.nsmap):
            for item in result.findall(".//Product", tree.nsmap):
                try:
                    asin = item.find(".//ASIN", tree.nsmap).text
                    unit = int(item.find(".//{*}PackageQuantity").text)
                    jan = result.attrib.get('Id')
                    title = item.find(".//{*}Title").text
                except AttributeError as e:
                    logger.debug(e)
                    continue

                cost = products_dict.get(jan)
                mws = MWS(asin=asin, title=title, jan=jan, unit=unit, cost=cost)
                mws_object_list.append(mws)

        return mws_object_list

    def get_competitive_pricing_for_asin(self, products: list, interval_sec: int = 2) -> list:
        logger.info('action=get_competitive_pricing_for_asin status=run')
        asin_price_dict = {}

        data_dict = dict(self.data)
        data_dict['MarketplaceId'] = settings.MARKETPLACEID
        data_dict['IdType'] = 'ASIN'
        data_dict['Action'] = 'GetCompetitivePricingForASIN'
        for index, asin in enumerate(products):
            data_dict[f'ASINList.ASIN.{str(index+1)}'] = asin

        url = self.create_request_url(data_dict=data_dict)
        response = request_api(url)
        time.sleep(interval_sec)
        tree = etree.fromstring(response.text)

        for item in tree.findall('.//{*}Product'):
            try:
                asin = item.find('.//{*}ASIN').text
            except AttributeError as ex:
                logger.error(ex)
                logger.error('asin is None')
                continue
            try:
                price = int(float(item.find('.//LandedPrice//Amount', tree.nsmap).text))
            except AttributeError as ex:
                logger.debug(ex)
                price = 0
            logger.debug(asin, price)
            asin_price_dict[asin] = price

        logger.info('action=get_competitive_pricing_for_asin status=done')

        return asin_price_dict

    def get_lowest_priced_offer_listtings_for_asin(self, products: list, interval_sec: int = 2) -> dict:
        logger.info('action=get_lowest_priced_offer_listtings_for_asin status=run')
        asin_price_dict = {}

        data_dict = dict(self.data)
        data_dict['MarketplaceId'] = settings.MARKETPLACEID
        data_dict['ItemCondition'] = 'New'
        data_dict['Action'] = 'GetLowestOfferListingsForASIN'
        for index, asin in enumerate(products):
            data_dict[f'ASINList.ASIN.{str(index+1)}'] = asin

        url = self.create_request_url(data_dict=data_dict)
        response = request_api(url)
        time.sleep(interval_sec)
        tree = etree.fromstring(response.text)

        for item in tree.findall('.//{*}Product'):
            try:
                asin = item.find('.//{*}ASIN').text
                price = int(float(item.find('.//LandedPrice//Amount', tree.nsmap).text))
            except AttributeError as e:
                logger.debug(e)
                continue
            asin_price_dict[asin] = price

        logger.info('action=get_competitive_pricing_for_asin status=done')

        return asin_price_dict

    def get_lowest_priced_offers_for_asin(self, asin: str, interval_sec: int = 1) -> int:
        logger.info('action=get_lowest_priced_offers_for_asin status=run')

        data_dict = dict(self.data)
        data_dict['MarketplaceId'] = settings.MARKETPLACEID
        data_dict['ItemCondition'] = 'New'
        data_dict['Action'] = 'GetLowestPricedOffersForASIN'
        data_dict['ASIN'] = asin

        query = urllib.parse.parse_qs(urllib.parse.urlparse(self.create_request_url(data_dict=data_dict)).query)
        url = "https://mws.amazonservices.jp/Products/2011-10-01"

        response = request_api(url=url, data=query)
        time.sleep(interval_sec)

        soup = BeautifulSoup(response.text, 'lxml')
        try:
            price = int(float(soup.select_one('LandedPrice Amount').text))
        except AttributeError as ex:
            logger.error(f'error={ex}')
            return None

        logger.info('action=get_lowest_priced_offers_for_asin status=done')
        return price

    def get_fee_my_fees_estimate(self, products: list, interval_sec: float = 0.5) -> list:
        logger.info('action=get_fee_my_fees_estimate status=run')
        asin_fees_list = []

        data_dict = dict(self.data)
        data_dict['Action'] = 'GetMyFeesEstimate'
        params = 'FeesEstimateRequestList.FeesEstimateRequest'
        for index, asin in enumerate(products):
            idx = str(index+1)
            data_dict[f'{params}.{idx}.IdType'] = 'ASIN'
            data_dict[f'{params}.{idx}.IdValue'] = asin
            data_dict[f'{params}.{idx}.Identifier'] = f'{idx}'
            data_dict[f'{params}.{idx}.IsAmazonFulfilled'] = 'true'
            data_dict[f'{params}.{idx}.MarketplaceId'] = 'A1VC38T7YXB528'
            data_dict[f'{params}.{idx}.PriceToEstimateFees.ListingPrice.Amount'] = '10000'
            data_dict[f'{params}.{idx}.PriceToEstimateFees.ListingPrice.CurrencyCode'] = 'JPY'

        url = self.create_request_url(data_dict=data_dict)
        response = request_api(url)
        time.sleep(interval_sec)

        try:
            tree = etree.fromstring(response.text)
        except etree.XMLSyntaxError as e:
            logger.error(e)
            logger.error(response.text)
            raise Exception

        for item in tree.findall(".//FeesEstimateResult", tree.nsmap):
            try:
                fee_rate = float(item.find(".//FeeDetailList//Amount", tree.nsmap).text) / 10000
                ship_fee = int(float(item.find('.//IncludedFeeDetailList//Amount', tree.nsmap).text))
            except AttributeError as e:
                logger.debug(e)
                fee_rate = 0.1
                ship_fee = 500
            asin = item.find('.//FeesEstimateIdentifier//IdValue', tree.nsmap).text
            asin_fees_list.append((asin, fee_rate, ship_fee))

        logger.info('action=get_fee_my_fees_estimate status=done')
        return asin_fees_list
    
    def pool_get_matching_product_for_id(self, products_list: list, filename: str) -> None:
        logger.info('action=pool_get_matching_product_for_id status=run')

        while products_list:
            products_five = [products_list.pop() for _ in range(5) if products_list]
            products_five_dict = {jan: cost for jan, cost in products_five}

            mws_objects_list = self.get_matching_product_for_id(products_five_dict)
            for mws_object in mws_objects_list:
                mws_object.filename = filename
                mws_object.save()
    
    def pool_get_competitive_pricing_for_asin(self, asin_list: list) -> None:
        logger.info('action=pool_get_competitive_pricing_for_asin status=run')

        while asin_list:
            asin_list_twenty = [asin_list.pop() for _ in range(20) if asin_list]
            asin_price_dict = self.get_competitive_pricing_for_asin(asin_list_twenty)
            for asin, price in asin_price_dict.items():
                MWS.update_price(asin=asin, price=price)

        logger.info('action=pool_get_competitive_pricing_for_asin status=done')
            
    def pool_get_lowest_priced_offers_for_asin(self):
        logger.info('action=pool_get_lowest_priced_offers_for_asin status=run')

        mws_products_list = MWS.get_price_is_None_products()
        for product in mws_products_list:
            price = self.get_lowest_priced_offers_for_asin(product.asin)
            if price is not None:
                MWS.update_price(asin=product.asin, filename=product.filename, price=price)

        logger.info('action=pool_get_lowest_priced_offers_for_asin status=done')

    def pool_get_fee_my_fees_estimate(self, asin_list: list):
        logger.info('action=pool_get_fee_my_fees_estimate status=run')

        while asin_list:
            asin_list_five = [asin_list.pop() for _ in range(5) if asin_list]
            asin_fees_list = self.get_fee_my_fees_estimate(asin_list_five)
            for asin, fee_rate, ship_fee in asin_fees_list:
                MWS.update_fee(asin=asin, fee_rate=fee_rate, shipping_fee=ship_fee)

    def request_report(self, report_type: str, start_date, end_date):
        logger.info('action=request_report status=run')
        xmlns = '{http://mws.amazonaws.com/doc/2009-01-01/}'
        param = dict(self.data)
        param['ReportType'] = report_type
        param['StartDate'] = start_date
        param['EndDate'] = end_date
        param['Action'] = 'RequestReport'
        param['Version'] = '2009-01-01'
        param['Merchant'] = param['SellerId']
        param.pop('SellerId')
        self.endpoint = '/'
        response = self.request(data_dict=param)
        request_id = response[0][0].findtext(xmlns + 'ReportRequestId')
        time.sleep(180)
        logger.info('action=request_report status=done')
        return request_id

    def get_report_request_list(self, request_id: str):
        logger.info('action=get_report_request_list status=run')
        xmlns = '{http://mws.amazonaws.com/doc/2009-01-01/}'
        param = dict(self.data)
        param['Action'] = 'GetReportRequestList'
        param['Version'] = '2009-01-01'
        param['ReportRequestIdList.Id.1'] = request_id
        param['Merchant'] = param['SellerId']
        param.pop('SellerId')
        self.endpoint = '/'
        response = self.request(data_dict=param)
        report_id = response[0][-1].findtext(xmlns + 'GeneratedReportId')
        time.sleep(60)
        logger.info('action=get_report_request_list status=done')
        return report_id

    def get_report(self, report_id):
        logger.info('action=get_report status=run')
        param = dict(self.data)
        param['Action'] = 'GetReport'
        param['Version'] = '2009-01-01'
        param['ReportId'] = report_id
        param['Merchant'] = param['SellerId']
        param.pop('SellerId')
        self.endpoint = '/'
        url = self.create_request_url(data_dict=param)
        response = request_api(url)
        response.encoding = response.apparent_encoding
        rows = response.text.split('\n')
        result = [row.split('\t') for row in rows]
        logger.info('action=get_report status=done')
        return result


def get_file_path() -> pathlib.Path:
    try:
        path = next(pathlib.Path(settings.SCRAPE_SCHEDULE_SAVE_PATH).iterdir())
    except StopIteration:
        logger.info('scraping_schedule_path is None')
        try:    
            path = next(pathlib.Path(settings.SCRAPE_SAVE_PATH).iterdir())
        except StopIteration:
            logger.info('scraping_path is None')
            return None
    return path

def main():
    logger.info('action=main status=run')

    get_matching_asin_for_jan_process = Process(target=main_get_matching_asin_for_jan)
    get_lowest_price_process = Process(target=main_get_lowest_price)
    get_fees_process = Process(target=main_get_fees)

    get_matching_asin_for_jan_process.start()
    get_lowest_price_process.start()
    get_fees_process.start()

    get_matching_asin_for_jan_process.join()
    get_lowest_price_process.join()
    get_fees_process.join()

def main_get_matching_asin_for_jan(interval_sec: int = 60) -> None:
    logger.info('action=main_get_matching_asin_for_jan status=run')
    
    while True:
        filepath = get_file_path()
        if filepath:
            df = pd.read_excel(str(filepath), dtype={'JAN': str}).drop_duplicates()
            product_list = df.to_numpy().tolist()
            client = AmazonClient()
            client.pool_get_matching_product_for_id(products_list=product_list, filename=filepath.stem)
            shutil.move(str(filepath), os.path.join(settings.SCRAPE_DONE_SAVE_PATH, filepath.name))
        else:
            time.sleep(interval_sec)


def main_get_lowest_price(interval_sec: int = 60) -> None:
    logger.info('action=main_get_lowest_price status=run')

    while True:
        asin_list = MWS.get_price_is_None_products()
        if asin_list:
            asin_list = list(map(lambda x: x[0], asin_list))
            client = AmazonClient()
            client.pool_get_competitive_pricing_for_asin(asin_list)
        else:
            time.sleep(interval_sec)


def main_get_fees(interval_sec: int = 60) -> None:
    logger.info('action=main_get_fees status=run')

    while True:
        asin_list = MWS.get_fee_is_None_products()
        if asin_list:
            asin_list = list(map(lambda x: x[0], asin_list))
            client = AmazonClient()
            client.pool_get_fee_my_fees_estimate(asin_list)
        else:
            time.sleep(interval_sec)
