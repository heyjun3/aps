from collections import defaultdict
import datetime

from flask import Flask
from flask import render_template
import pandas as pd

from keepa.models import KeepaProducts
from keepa import keepa 


app = Flask(__name__)


def convert_price_rank_data(price_data, rank_data):

    rank_dict = {keepa.convert_keepa_time_to_datetime_date(int(k)): v for k, v in rank_data.items()}
    price_dict = {keepa.convert_keepa_time_to_datetime_date(int(k)): v for k, v in price_data.items()}

    rank_df = pd.DataFrame(data=list(rank_dict.items()), columns=['date', 'rank']).astype({'rank': int})
    price_df = pd.DataFrame(data=list(price_dict.items()), columns=['date', 'price']).astype({'price': int})

    df = pd.merge(rank_df, price_df, on='date', how='outer')
    df = df.fillna(method='ffill')
    delay = datetime.datetime.now().date() - datetime.timedelta(days=90)
    df = df[df['date'] > delay]
    df = df.sort_values('date', ascending=True)
    products = df.to_dict('records')
    
    return products


@app.route('/')
def hello():
    asin = 'B08F59Z1B8'
    asin_1 = 'B08L3HDFST'
    product_1 = KeepaProducts.get_keepa_product(asin)
    product_2 = KeepaProducts.get_keepa_product(asin_1)

    p_1 = convert_price_rank_data(product_1.price_data, product_1.rank_data)
    p_2 = convert_price_rank_data(product_2.price_data, product_2.rank_data)
    product_list = []
    product_list.append({'product': p_1, 'asin': asin})
    product_list.append({'product': p_2, 'asin': asin_1})


    return render_template('chart.html', products=product_list)


def start():
    app.run(host='127.0.0.1', port='8080', threaded=True, debug=True)
