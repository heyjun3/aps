from datetime import datetime
import urllib.parse
import hmac
import six
import hashlib
import base64
import time
from multiprocessing import Queue
from multiprocessing import Manager
from multiprocessing import Process
import logging.config
import os
import pathlib
import shutil

import openpyxl
import requests
import xml.etree.ElementTree as et
import pandas as pd
from lxml import etree

from mws import multiprocess
from mws.models import MWS
import settings


logger = logging.getLogger(__name__)


def datetime_encode(dt):
    return dt.strftime('%Y-%m-%dT%H:%M:%SZ')


def request_api(url):
    logger.info('action=request_api status=run')

    for _ in range(60):
        try:
            response = requests.post(url, timeout=30.0)
            if not response.status_code == 200 or response is None:
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
        # print(response_text)
        response_xml = et.fromstring(response_text)
        # xml = et.fromstring(response_text)
        # xml = parse(response_text)
        # xml = xml.getroot()

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

    def get_matching_product_for_id(self, products_dict: dict, filename: str, price_que=None, fee_que=None, manager=None):
        """Searching Asin code for Jan code
        share memory list append searching object
        share.append([jan, cost, asin, rank, quantity])"""
        logger.info('action=get_matching_product_for_id status=run')
        data = []
        products = list(products_dict.keys())

        while products:
            product = [products.pop() for _ in range(5) if products]

            data_dict = dict(self.data)
            data_dict['MarketplaceId'] = settings.MARKETPLACEID
            data_dict['IdType'] = 'JAN'
            data_dict['Action'] = 'GetMatchingProductForId'
            for index, jan in enumerate(product):
                data_dict[f'IdList.Id.{str(index+1)}'] = str(jan)
            logger.debug(data_dict)

            url = self.create_request_url(data_dict=data_dict)
            response = request_api(url)
            response.encoding = response.apparent_encoding
            logger.info('action=get_matching_product_for_id status=run')
            # import xml.dom.minidom
            # x = xml.dom.minidom.parseString(response.text)
            # print(x.toprettyxml())
            time.sleep(1)
            tree = etree.fromstring(response.text)

            asin_lst = []
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
                    logger.debug(asin, unit, jan)
                    data.append([asin, unit, jan])
                    asin_lst.append(asin)
                    cost = products_dict.get(jan)
                    mws = MWS(asin=asin, title=title, jan=jan, unit=unit, filename=filename, cost=cost)
                    mws.save()

            if price_que is not None and fee_que is not None:
                price_que.put(asin_lst)
                fee_que.put(asin_lst)

        df = pd.DataFrame(data=data, columns=['asin', 'unit', 'jan']).astype({'unit': int})
        logger.info('action=get_matching_product_for_id status=done')
        manager['matching_df'] = df


    def get_competitive_pricing_for_asin(self, products, filename: str):
        logger.info('action=get_competitive_pricing_for_asin status=run')
        data = []

        while products:
            product = [products.pop() for _ in range(20) if products]
            data_dict = dict(self.data)
            data_dict['MarketplaceId'] = settings.MARKETPLACEID
            data_dict['IdType'] = 'ASIN'
            data_dict['Action'] = 'GetCompetitivePricingForASIN'
            for index, asin in enumerate(product):
                data_dict[f'ASINList.ASIN.{str(index+1)}'] = asin

            url = self.create_request_url(data_dict=data_dict)
            response = request_api(url)
            time.sleep(0.1 * len(product))
            tree = etree.fromstring(response.text)

            for item in tree.findall('.//{*}Product'):
                try:
                    asin = item.find('.//{*}ASIN').text
                    price = int(float(item.find('.//LandedPrice//Amount', tree.nsmap).text))
                except AttributeError as e:
                    logger.debug(e)
                    continue
                logger.debug(asin, price)
                data.append([asin, price])
                MWS.update_price(asin=asin, filename=filename, price=price)

        logger.info('action=get_competitive_pricing_for_asin status=done')
        return data

    def get_lowest_priced_offer_listtings_for_asin(self, products: list, filename: str):
        logger.info('action=get_lowest_priced_offer_listtings_for_asin status=run')
        data = []

        while products:
            product = [products.pop() for _ in range(20) if products]
            data_dict = dict(self.data)
            data_dict['MarketplaceId'] = settings.MARKETPLACEID
            data_dict['ItemCondition'] = 'New'
            data_dict['Action'] = 'GetLowestOfferListingsForASIN'
            for index, asin in enumerate(product):
                data_dict[f'ASINList.ASIN.{str(index+1)}'] = asin

            url = self.create_request_url(data_dict=data_dict)
            response = request_api(url)
            time.sleep(0.1 * len(product))
            tree = etree.fromstring(response.text)

            for item in tree.findall('.//{*}Product'):
                try:
                    asin = item.find('.//{*}ASIN').text
                    price = int(float(item.find('.//LandedPrice//Amount', tree.nsmap).text))
                except AttributeError as e:
                    logger.debug(e)
                    continue
                logger.debug(asin, price)
                data.append([asin, price])
                MWS.update_price(asin=asin, filename=filename, price=price)

        logger.info('action=get_competitive_pricing_for_asin status=done')
        return data

    def get_fee_my_fees_estimate(self, products, filename: str):
        logger.info('action=get_fee_my_fees_estimate status=run')
        data = []

        while products:
            product = [products.pop() for _ in range(5) if products]
            data_dict = dict(self.data)
            data_dict['Action'] = 'GetMyFeesEstimate'
            params = 'FeesEstimateRequestList.FeesEstimateRequest'
            for index, asin in enumerate(product):
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
            time.sleep(0.1 * len(product))

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
                logger.debug(asin, fee_rate, ship_fee)
                data.append([asin, fee_rate, ship_fee])
                MWS.update_fee(asin=asin, filename=filename, fee_rate=fee_rate, shipping_fee=ship_fee)

        logger.info('action=get_fee_my_fees_estimate status=done')
        return data

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


def open_excel_file(filepath: str) -> dict:
    """Open excel file Return row list object"""
    workbook = openpyxl.load_workbook(filepath)
    worksheet = workbook[workbook.sheetnames[0]]
    worksheet.delete_rows(1)
    # values = list(set(list(worksheet.values)))
    # return list(map(list, values))
    return dict(sorted(list(worksheet.values), key=lambda x: x[1], reverse=True))


def get_file_path():
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
    while True:
        file = get_file_path()
        if file:
            client = AmazonClient()
            products_df = pd.read_excel(str(file), dtype={'JAN': str}).rename(columns={'JAN': 'jan', 'Cost': 'cost'}).drop_duplicates()
            product_dict = {jan: cost for jan, cost in zip(products_df['jan'], products_df['cost'])}
            price_que = Queue()
            fee_que = Queue()
            manager = Manager()
            manager = manager.dict()
            filename = file.stem

            get_matching_prodcut_for_id_process = Process(target=client.get_matching_product_for_id, args=(product_dict, filename, price_que, fee_que, manager))
            get_competitive_pricing_for_asin_process = Process(target=multiprocess.get_competitive_pricing_for_asin_worker, args=(filename, price_que, manager))
            get_fee_my_fees_estimate_process = Process(target=multiprocess.get_fee_my_fees_estimate_worker, args=(filename, fee_que, manager))

            get_matching_prodcut_for_id_process.start()
            get_competitive_pricing_for_asin_process.start()
            get_fee_my_fees_estimate_process.start()

            get_matching_prodcut_for_id_process.join()
            price_que.put(None)
            fee_que.put(None)

            get_competitive_pricing_for_asin_process.join()
            get_fee_my_fees_estimate_process.join()

            matching_df = manager.get('matching_df')
            price_df = manager.get('price_df')
            fees_df = manager.get('fees_df')

            df = matching_df.merge(products_df, on='jan', how='inner').merge(price_df, on='asin', how='inner')\
                .merge(fees_df, on='asin', how='inner')

            df['profit'] = df['price'] - (df['cost'] * df['unit']) - ((df['price'] * df['fee_rate']) * 1.1) - df['ship_fee']
            df['profit_rate'] = df['profit'] / df['price']

            df = df.query('profit >= 200 and profit_rate >= 0.1').astype({'profit': int})
            df.to_pickle(os.path.join(settings.MWS_SAVE_PATH, f'{file.stem}.pickle'))
            shutil.move(str(file), os.path.join(settings.SCRAPE_DONE_SAVE_PATH, file.name))
        else:
            time.sleep(60)
