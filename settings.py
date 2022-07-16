import os
import configparser


BASE_PATH = os.path.dirname(__file__)
config = configparser.ConfigParser()
config.read(os.path.join(BASE_PATH, 'settings.ini'))
LOGGING_CONF_PATH = os.path.join(BASE_PATH, 'logging.conf')

# DATABASE settings
USERNAME = config['DB']['USERNAME']
PASSWORD = config['DB']['PASSWORD']
HOST = config['DB']['HOST']
PORT = config['DB']['PORT']
DB = config['DB']['DB']
TEST_DB = config['DB']['TEST']
DB_URL = f"postgresql+psycopg2://{USERNAME}:{PASSWORD}@{HOST}:{PORT}/{DB}"
DB_ASYNC_URL = f"postgresql+asyncpg://{USERNAME}:{PASSWORD}@{HOST}:{PORT}/{DB}"
DB_TEST_URL = f"postgresql+asyncpg://{USERNAME}:{PASSWORD}@{HOST}:{PORT}/{TEST_DB}"

# REDIS settings
REDIS_HOST = config['REDIS']['HOST']
REDIS_PORT = config['REDIS']['PORT']
REDIS_DB = config['REDIS']['DB']

# WEB settings
HOST = config['WEB']['HOST']
PORT = config['WEB']['PORT']

# SAVE PATH
BASE_SAVE_PATH = os.path.join(BASE_PATH, 'scrape_files')
SCRAPE_SAVE_PATH = os.path.join(BASE_SAVE_PATH, 'scrape')
SCRAPE_SCHEDULE_SAVE_PATH =  os.path.join(BASE_SAVE_PATH, 'scrape_schedule')
SCRAPE_DONE_SAVE_PATH = os.path.join(BASE_SAVE_PATH, '.scrape')
MWS_SAVE_PATH = os.path.join(BASE_SAVE_PATH, 'mws')
MWS_DONE_SAVE_PATH = os.path.join(BASE_SAVE_PATH, '.mws')
KEEPA_SAVE_PATH = os.path.join(BASE_SAVE_PATH, 'keepa')
SOURCE_PATH = os.path.join(BASE_PATH, 'ims', 'source.xlsx')

# MWS settings
SECRET_KEY = config['MWS']['SECRET_KEY']
MARKETPLACEID = config['MWS']['MARKETPLACEID']
ACCESS_ID = config['MWS']['ACCESS_ID']
SELLER_ID = config['MWS']['SELLER_ID']
DOMAIN = 'mws.amazonservices.jp'
ENDPOINT = '/Products/2011-10-01'
XMLNS = '{http://mws.amazonservices.com/schema/Products/2011-10-01}'
NS2 = '{http://mws.amazonservices.com/schema/Products/2011-10-01/default.xsd}'

# SPAPI settings
REFRESH_TOKEN = config['SPAPI']['REFRESH_TOKEN']
CLIENT_ID = config['SPAPI']['CLIENT_ID']
CLIENT_SECRET = config['SPAPI']['CLIENT_SECRET']
AWS_ACCESS_ID = config['SPAPI']['AWS_ACCESS_ID']
AWS_SECRET_KEY = config['SPAPI']['AWS_SECRET_KEY']

# BUFFALO settings
BUFFALO_URL = 'https://www.buffalo-direct.com/'
BUFFALO_START_URL = 'https://www.buffalo-direct.com/directshop/products/list_category.php?category_id=1181'
BUFFALO_DETAIL_PAGE_URL = 'https://www.buffalo-direct.com/directshop/products/detail.php'

# KEEPA settings
KEEPA_ACCESS_KEY = config['KEEPA']['KEEPA_ACCESS_KEY']

# NETSEA settings
NETSEA_LOGIN_URL = 'https://www.netsea.jp/login'
NETSEA_NEXT_URL = 'https://www.netsea.jp/search/'
NETSEA_SHOP_URL = 'https://www.netsea.jp/shop/'
NETSEA_ENDPOINT = 'https://www.netsea.jp'
NETSEA_ID = config['NETSEA']['ID']
NETSEA_PASSWD = config['NETSEA']['PASSWD']

# RAKUTEN settings
RAKUTEN_APP_ID = config['RAKUTEN']['APP_ID']
REQUEST_URL = 'https://app.rakuten.co.jp/services/api/IchibaItem/Search/20170706'
RAKUTEN_ENDPOINT = 'https://item.rakuten.co.jp/'

# SUPER settings
SUPER_LOGIN_URL = 'https://www.superdelivery.com/p/login.do'
SUPER_LOGIN_INFO = {
            'identification': config['SUPER']['MAIL_ADDRESS'],
            'password': config['SUPER']['PASSWORD'],
}
SUPER_DOMAIN_URL = 'https://www.superdelivery.com'
SUPER_NEW_PRODUCTS_URL = 'https://www.superdelivery.com/p/do/psl/'

# PC4U settings
PC4U_ENDPOINT = 'https://www.pc4u.co.jp/'