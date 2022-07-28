import math

from flask import Flask
from flask import request
from flask import jsonify
from flask_cors import CORS
from keepa.keepa import keepa_request_products, scrape_keepa_request
from keepa.models import KeepaProducts
from mws.api import AmazonClient

from mws.models import MWS
from spapi.models import AsinsInfo
import settings
import log_settings


logger = log_settings.get_logger(__name__)
app = Flask(__name__)
CORS(app, supports_credentials=True, origins=['http://192.168.0.8:5000'])


@app.route('/list', methods=['GET'])
def get_list():
    if request.method == 'GET':
        filename_list = MWS.get_filenames()
        return jsonify({'list': filename_list}), 200
    else:
        return jsonify({'error': 'Bad request method'}), 404


@app.route('/deleteFile/<string:filename>', methods=['DELETE'])
def delete_file(filename):
    if request.method == 'DELETE':
        flag = MWS.delete_rows(filename)
        logger.info(filename)
        if flag:
            return jsonify(None), 200
        return jsonify({'action': 'delete_file', 'status': 'error'}), 400


@app.route('/search/<string:asin>', methods=['GET'])
def chart_render(asin: str):
    if request.method == 'GET':
        asin = asin.strip()
        product = KeepaProducts.get_keepa_product(asin)
        if not product:
            response = keepa_request_products([asin])
            asin, drops, price_data, rank_data = scrape_keepa_request(response)[0]
            KeepaProducts.update_or_insert(asin, drops, price_data, rank_data)

        title = AsinsInfo.get_title(asin)
        if not title:
            client = AmazonClient()
            title = client.get_matching_product_for_asin(asin)
        product.render_data['title'] = title
        return jsonify(product.render_data), 200
    else:
        return jsonify({'status': 'error'}), 400


@app.route('/chart_list/<string:filename>', methods=['GET'])
def get_chart_data(filename: str) -> str:
    if request.method == 'GET':
        chart_data = []
        current_page_number = int(float(request.args.get('page', '1')))
        count = int(float(request.args.get('count', '100')))

        products = MWS.get_chart_data(filename=filename, page=current_page_number, count=count)
        max_page = math.ceil(MWS.get_row_count(filename) / count)
        if not products:
            return jsonify({'status': 'error', 'message': 'chart data is None'}), 200
        for mws, data in products:
            data['asin'] = mws.asin
            data['jan'] = mws.jan
            data['title'] = mws.title
            chart_data.append(data)

        return jsonify({'chart_data': chart_data, 'current_page': current_page_number, 'max_page': max_page}), 200
    else:
        return jsonify({'status': 'error'}), 400


def start():
    app.run(host=settings.HOST, port=settings.PORT, threaded=True, debug=True)
