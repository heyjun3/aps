from flask import Flask
from flask import render_template
from flask import redirect
from flask import url_for
from flask import request

from mws.models import MWS
import settings


app = Flask(__name__)


@app.route('/')
def index():
    filename_list = MWS.get_completion_filename_list()

    return render_template('index.html', save_path=filename_list)


@app.route('/graph/<string:filename>', methods=['GET'])
def view_graph(filename):
    products_list = MWS.get_render_data(filename=filename)
    return render_template('chart.html', products=products_list)


@app.route('/delete/<string:filename>', methods=['POST'])
def delete_filename(filename):
    MWS.delete_objects(filename)
    return redirect(url_for('index'))


@app.route('/chart/<string:filename>', methods=['GET'])
def chart(filename):
    render_data_list = []
    products_list = MWS.get_render_data(filename=filename)
    for mws, keepa in products_list:
        date, rank, price = keepa.render_price_rank_data_list
        date = list(map(lambda x: x.isoformat(), date))
        product = {
            'title': mws.title,
            'asin': mws.asin,
            'date': date,
            'price': price,
            'rank': rank,
            'jan': mws.jan,
        }
        render_data_list.append(product)

    return render_template('chart_js.html', products_list=render_data_list)


def start():
    app.run(host=settings.HOST, port=settings.PORT, threaded=True, debug=True)
