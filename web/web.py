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


@app.route('/delete', methods=['POST'])
def delete_filename():
    filename_list = request.form.getlist('filename')
    if not filename_list:
        return redirect(url_for('index'))
    MWS.delete_objects(filename_list)
    return redirect(url_for('index'))


@app.route('/chart/<string:filename>', methods=['GET'])
def chart(filename):
    render_data_list = []
    products_list = MWS.get_render_data(filename=filename)
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

    return render_template('chart.html', products_list=render_data_list)


def start():
    app.run(host=settings.HOST, port=settings.PORT, threaded=True, debug=True)
