import urllib.parse
import time

from crawler.netsea.netsea import Netsea
from crawler.netsea.netsea import NetseaHTMLPage
from crawler import utils
import settings
import log_settings


logger = log_settings.get_logger(__name__)


def run_netsea_at_shop_id(shop_id: str, path: str = 'search') -> None:
    url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)
    params = {'sort': 'PD', 'supplier_id': shop_id, 'ex_so': 'Y', 'searched': 'Y'}
    client = Netsea(url, params)
    client.start_search_products()


def run_get_all_shop_info(path: str = 'shop', interval_sec: int = 2) -> None:
    logger.info('action=run_get_all_shop_info status=run')
    url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)

    for index in range(1, 9):
        params = {'category_id': str(index), 'sort': 'NEW'}
        response = utils.request(url=url, params=params)
        shops = NetseaHTMLPage.scrape_shop_list_page(response.text)
        list(map(lambda x: x.save(), shops))
        time.sleep(interval_sec)