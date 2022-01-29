import datetime
import pathlib
import os

from flask import Flask
from flask import render_template
import pandas as pd

from keepa.models import KeepaProducts
from keepa import keepa 
import settings


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
def index():
    path = list(pathlib.Path(settings.KEEPA_SAVE_PATH).iterdir())
    path = list(map(lambda x: x.stem, path))
    print(path)

    return render_template('index.html', save_path=path)


@app.route('/graph/<string:filename>', methods=['GET'])
def view_graph(filename):
    path = os.path.join(settings.KEEPA_SAVE_PATH, f'{filename}.xlsx')
    df = pd.read_excel(path)
    asin_list = list(df['asin'])
    print(len(asin_list))
    products_list = []
    for asin in asin_list:
        keepa_product = KeepaProducts.get_keepa_product(asin)
        if keepa_product is None or keepa_product.price_data is None or keepa_product.rank_data is None:
            print(asin)
            continue
        price_rank_data = convert_price_rank_data(keepa_product.price_data, keepa_product.rank_data)
        products_list.append({'product': price_rank_data, 'asin': asin})
    print(len(products_list))
    return render_template('chart.html', products=products_list)


def start():
    app.run(host=settings.HOST, port='8080', threaded=True, debug=True)



