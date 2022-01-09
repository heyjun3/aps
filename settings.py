import os
import configparser


BASE_PATH = os.path.dirname(__file__)
config = configparser.ConfigParser()
config.read(os.path.join(BASE_PATH, 'settings.ini'))

LOGGING_CONF_PATH = os.path.join(BASE_PATH, 'logging.conf')

SECRET_KEY = config['MWS']['SECRET_KEY']
MARKETPLACEID = config['MWS']['MARKETPLACEID']
ACCESS_ID = config['MWS']['ACCESS_ID']
SELLER_ID = config['MWS']['SELLER_ID']
DOMAIN = 'mws.amazonservices.jp'
ENDPOINT = '/Products/2011-10-01'
XMLNS = '{http://mws.amazonservices.com/schema/Products/2011-10-01}'
NS2 = '{http://mws.amazonservices.com/schema/Products/2011-10-01/default.xsd}'

KEEPA_ACCESS_KEY = config['KEEPA']['KEEPA_ACCESS_KEY']
USERNAME = config['DB']['USERNAME']
PASSWORD = config['DB']['PASSWORD']
HOST = config['DB']['HOST']
PORT = config['DB']['PORT']
DB = config['DB']['DB']
DB_URL = f"postgresql://{USERNAME}:{PASSWORD}@{HOST}:{PORT}/{DB}"


BASE_SAVE_PATH = os.path.join(BASE_PATH, 'scrape_files')

SCRAPE_SAVE_PATH = os.path.join(BASE_SAVE_PATH, 'scrape')
SCRAPE_SCHEDULE_SAVE_PATH =  os.path.join(BASE_SAVE_PATH, 'scrape_schedule')
SCRAPE_DONE_SAVE_PATH = os.path.join(BASE_SAVE_PATH, '.scrape')
MWS_SAVE_PATH = os.path.join(BASE_SAVE_PATH, 'mws')
MWS_DONE_SAVE_PATH = os.path.join(BASE_SAVE_PATH, '.mws')
KEEPA_SAVE_PATH = os.path.join(BASE_SAVE_PATH, 'keepa')
SOURCE_PATH = os.path.join(BASE_PATH, 'ims', 'source.xlsx')

BUFFALO_URL = 'https://www.buffalo-direct.com/'
BUFFALO_START_URL = 'https://www.buffalo-direct.com/directshop/products/list_category.php?category_id=1181'