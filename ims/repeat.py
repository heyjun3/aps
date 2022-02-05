import datetime
import os

import pandas as pd

from ims.models import Product
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

    df = pd.DataFrame(data=None, columns={'sku': str, 'asin': str})
    for data in inventory_data[1:-1]:
        sku, asin = data[2], data[11]
        df = df.append({'sku': sku, 'asin': asin}, ignore_index=True)

    return df

def main():

    netsea_df = netsea.collect_favorite_products()
    super_df = super.collection_favorite_products()
    mws_df = get_marchant_listings_inactive_data()
    products = Product.get_all_objects()
    products = list(map(lambda x: x.value, products))
    product_df = pd.DataFrame(data=products)

    df_concat = pd.concat([netsea_df, super_df])
    df = pd.merge(mws_df, product_df, on='asin', how='inner')
    df = pd.merge(df, df_concat, on='jan', how='inner')
    df = df[['jan', 'cost']].dropna().rename(columns={'jan': 'JAN', 'cost': 'Cost'})

    timestamp = datetime.datetime.now().strftime('%Y%m%d_%H%M%S')
    save_path = os.path.join(settings.SCRAPE_SCHEDULE_SAVE_PATH, f'repeatedly{timestamp}.xlsx')
    df.to_excel(save_path, index=False)
