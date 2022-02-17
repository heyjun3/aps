import urllib.parse

from crawler.netsea.netsea import Netsea
import settings


def run_netsea_at_shop_id(shop_id: str, path: str = 'search') -> None:
    url = urllib.parse.urljoin(settings.NETSEA_ENDPOINT, path)
    params = {'sort': 'PD', 'supplier_id': shop_id, 'ex_so': 'Y', 'searched': 'Y'}
    client = Netsea(url, params)
    client.start_search_products()