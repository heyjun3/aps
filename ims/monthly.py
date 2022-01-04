import os

import openpyxl
import datetime
from dateutil.relativedelta import relativedelta

from ims.models import Product, Stock
from mws.api import AmazonClient
from ims import models

PATH = os.path.join(os.path.dirname(__file__), 'source.xlsx')

def insert():
    workbook = openpyxl.open(PATH)
    sheet = workbook['insert']
    sheet.delete_rows(1)

    for row in sheet.values:
        date, name, asin, jan, sku, fnsku, danger_class, sell_price, cost_price = row
        Product.create(date, name, asin, jan, sku, fnsku, danger_class, sell_price, cost_price)

    workbook.close()


def purchasing_data_update():
    workbook = openpyxl.open(PATH)
    sheet = workbook['purchasing']
    sheet.delete_rows(1)

    for row in sheet.values:
        date, info, asin, sku, stock, unit_price, *_ = row
        # print(date, info, asin, sku, stock, unit_price)
        Product.update_cost_price(sku, int(unit_price))
        Stock.add_home_stock(sku=sku, home_stock=int(stock))

    workbook.close()


def get_fba_fulfillment_inventory_receipts_data():
    report_type = '_GET_FBA_FULFILLMENT_INVENTORY_RECEIPTS_DATA_'
    tz_jst = datetime.timezone(datetime.timedelta(hours=9))
    today = datetime.date.today()
    end = datetime.datetime(today.year, today.month, day=1, tzinfo=tz_jst)
    start = end - relativedelta(months=1)
    start = start.isoformat()
    end = end.isoformat()

    amazon_client = AmazonClient()
    request_id = amazon_client.request_report(report_type=report_type, start_date=start, end_date=end)
    report_id = amazon_client.get_report_request_list(request_id=request_id)
    inventory_data = amazon_client.get_report(report_id)

    for row in inventory_data[1:]:
        date, fnsku, sku, product_name, quantity, *_ = row
        Stock.decrease_home_stock(sku=sku, home_stock=int(quantity))


def get_fba_fulfillment_current_inventory_data():
    report_type = '_GET_FBA_FULFILLMENT_CURRENT_INVENTORY_DATA_'
    tz_jst = datetime.timezone(datetime.timedelta(hours=9))
    today = datetime.date.today()
    end = datetime.datetime(today.year, today.month, day=1, tzinfo=tz_jst)
    start = end - datetime.timedelta(days=1)
    start = start.isoformat()
    end = end.isoformat()

    amazon_client = AmazonClient()
    request_id = amazon_client.request_report(report_type=report_type, start_date=start, end_date=end)
    report_id = amazon_client.get_report_request_list(request_id=request_id)
    inventory_data = amazon_client.get_report(report_id)

    Stock.initialize_fba_stock()
    for row in inventory_data[1:]:
        data, fnsku, sku, product_name, quantity, *_ = row
        Stock.add_fba_stock(sku=sku, fba_stock=int(quantity))


def main():
    insert()  # insert amazonproducts database
    purchasing_data_update()  # update cost_price and add home stocks
    get_fba_fulfillment_inventory_receipts_data()  # get inventory receipts and decrease home stocks
    get_fba_fulfillment_current_inventory_data()  # get end of last month inventory
    models.calc_stock_cost()

