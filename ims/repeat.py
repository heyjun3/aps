import datetime
import os

import openpyxl

from ims.models import InactiveStock
from ims.models import FavoriteProduct
from mws.api import AmazonClient
from crawler.netsea import netsea
from crawler.super import super
import settings


def get_marchant_listings_inactive_data():
    report_type = '_GET_MERCHANT_LISTINGS_INACTIVE_DATA_'
    tz_jst = datetime.timezone(datetime.timedelta(hours=9))
    today = datetime.date.today()
    end = datetime.datetime(today.year, today.month, today.day, tzinfo=tz_jst)
    start = end - datetime.timedelta(days=1)
    start = start.isoformat()
    end = end.isoformat()

    amazon_client = AmazonClient()
    request_id = amazon_client.request_report(report_type=report_type, start_date=start, end_date=end)
    report_id = amazon_client.get_report_request_list(request_id=request_id)
    inventory_data = amazon_client.get_report(report_id)

    for data in inventory_data[1:-1]:
        sku, asin = data[2], data[11]
        InactiveStock.save(sku, asin)


def main():
    InactiveStock.delete()
    FavoriteProduct.delete()

    netsea.collect_favorite_products()
    super.collection_favorite_products()
    get_marchant_listings_inactive_data()

    data = {}
    products = InactiveStock.get_asin_cost()
    for product in products:
        inactive, master, favorite = product
        data[master.jan] = favorite.cost

    workbook = openpyxl.Workbook()
    sheet = workbook['Sheet']
    sheet.append(['JAN', 'Cost'])

    for key, value in data.items():
        sheet.append([key, value])

    timestamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')
    save_path = os.path.join(settings.SCRAPE_SCHEDULE_SAVE_PATH, f'repeatedly{timestamp}.xlsx')
    workbook.save(save_path)
    workbook.close()
