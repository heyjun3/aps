import math
import datetime

from flask import Flask
from flask import render_template
from flask import redirect
from flask import url_for
from flask import request
from flask import jsonify
from flask_cors import CORS
from keepa.models import KeepaProducts
from keepa.models import convert_keepa_time_to_datetime_date
import pandas as pd
import numpy as np

from mws.models import MWS
import settings
import log_settings


logger = log_settings.get_logger(__name__)
app = Flask(__name__)
CORS(
    app, supports_credentials=True,
)


@app.route('/')
def index():
    filename_list = MWS.get_done_filenames()

    return render_template('index.html', save_path=filename_list)


@app.route('/list', methods=['GET'])
def get_list():
    if request.method == 'GET':
        filename_list = MWS.get_done_filenames()
        return jsonify({'list': filename_list}), 200
    else:
        return jsonify({'error': 'Bad request method'}), 404


@app.route('/delete', methods=['POST'])
def delete_filename():
    filename_list = request.form.getlist('filename')
    if not filename_list:
        return redirect(url_for('index'))
    MWS.delete_objects(filename_list)
    return redirect(url_for('index'))


@app.route('/deleteFile/<string:filename>', methods=['DELETE'])
def delete_file(filename):
    if request.method == 'DELETE':
        flag = MWS.delete_rows(filename)
        logger.info(filename)
        if flag:
            return jsonify(None), 200
        return jsonify({'action': 'delete_file', 'status': 'error'}), 400


@app.route('/chart/<string:filename>', methods=['GET'])
def chart(filename):
    count = 500
    render_data_list = []
    
    current_page_num = int(request.args.get('page', '1'))
    products_list = MWS.get_render_data(filename=filename, page=current_page_num, count=count)
    max_pages = math.ceil(MWS.get_max_row_count(filename) / count)
    if not products_list:
        return redirect(url_for('index'))
        
    for mws, keepa in products_list:
        if keepa is None:
            continue
        elif keepa.get('date') and keepa.get('rank') and keepa.get('price'):
            product = {
                'title': mws.title,
                'asin': mws.asin,
                'date': keepa.get('date'),
                'price': keepa.get('price'),
                'rank': keepa.get('rank'),
                'jan': mws.jan,
            }
            render_data_list.append(product)
    logger.info({'action': 'chart', 'status': 'success'})
    return render_template('chart.html',
                           products_list=render_data_list,
                           current_page_num=current_page_num,
                           max_pages=max_pages,
                        )

@app.route('/search/<string:asin>', methods=['GET'])
def chart_render(asin: str):
    if request.method == 'GET':
        asin = asin.strip()
        product = KeepaProducts.get_keepa_product(asin)
        if product:
            df = pd.DataFrame(data=list(product.price_data.items()), columns=['date', 'price']).astype({'date':int,  'price': int})
            df = df.replace(-1.0, np.nan)
            df = df.fillna(method='ffill')
            df = df.fillna(method='bfill')
            df = df.replace([np.nan], [None])
            df['date'] = df['date'].map(convert_keepa_time_to_datetime_date)
            past_date = datetime.datetime.now().date() - datetime.timedelta(days=90)
            df = df[df['date'] > past_date]
            df = df.sort_values('date', ascending=True)
            df['date'] = df['date'].map(lambda x: x.isoformat())
            return jsonify(df.to_dict(orient='records')), 200
        else:
            return jsonify({'status': 'error'}), 400


def start():
    app.run(host=settings.HOST, port=settings.PORT, threaded=True, debug=True)
