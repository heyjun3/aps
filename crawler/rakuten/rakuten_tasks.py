from cmath import log
import log_settings
from crawler.rakuten.rakuten import RakutenAPIClient

logger = log_settings.get_logger(__name__)


SHOP_CODES = ['ksdenki', 'dj', 'e-zoa', 'reckb', 'jtus', 'ioplaza', 'ikebe']
SHOP_CODES = ['ioplaza']


def run_rakuten_search_at_shop_code(shop_code: str) -> None:
    logger.info('action=run_rakuten_search_at_shop_code status=run')

    rakuten = RakutenAPIClient(shop_code=shop_code)
    rakuten.run_rakuten_search()
    
    logger.info('action=run_rakuten_search_at_shop_code status=run')


def run_rakuten_search_all() -> None:
    logger.info('action=run_rakuten_search_all status=run')

    for shop_code in SHOP_CODES:
        run_rakuten_search_at_shop_code(shop_code)

    logger.info('action=run_rakuten_search_all status=done')
