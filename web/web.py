import math

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

    return render_template('chart.html',
                           products_list=render_data_list,
                           current_page_num=current_page_num,
                           max_pages=max_pages,
                        )


def start():
    app.run(host=settings.HOST, port=settings.PORT, threaded=True, debug=True)
