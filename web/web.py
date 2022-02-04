from flask import Flask
from flask import render_template

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


def start():
    app.run(host=settings.HOST, port=settings.PORT, threaded=True, debug=True)
