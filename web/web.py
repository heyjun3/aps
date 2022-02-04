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


@app.route('/')
def index():
    path = list(pathlib.Path(settings.KEEPA_SAVE_PATH).iterdir())
    path = list(map(lambda x: x.stem, path))
    path = sorted(path, key=lambda x: x)

    return render_template('index.html', save_path=path)


@app.route('/graph/<string:filename>', methods=['GET'])
def view_graph(filename):
    path = os.path.join(settings.KEEPA_SAVE_PATH, f'{filename}.xlsx')
    df = pd.read_excel(path)
    asin_list = list(df['asin'])
    products_list = []
    for asin in asin_list:
        keepa_product = KeepaProducts.get_keepa_product(asin)
        if keepa_product is None or keepa_product.price_data is None or keepa_product.rank_data is None:
            continue
        products_list.append({'product': keepa_product.render_price_rank_data, 'asin': asin})
    return render_template('chart.html', products=products_list)


def start():
    app.run(host=settings.HOST, port=settings.PORT, threaded=True, debug=True)
