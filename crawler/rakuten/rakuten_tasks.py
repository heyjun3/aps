from cmath import log
import log_settings
from crawler.rakuten.rakuten import RakutenAPIClient, RakutenCrawler

logger = log_settings.get_logger(__name__)


SHOP_CODES = [
    'superdeal', 'ksdenki', 'dj', 'e-zoa', 'reckb', 
    'ioplaza', 'ikebe',
    'shoptsukumo', 'premiumgt', 'acer-direct', 'wakeari', 'sakurayama',
    'ikeshibu', 'aikyoku-bargain-center', 'aikyoku', 'ikebe-rockhouse', 'ishibashi',
    'pckoubou', 'e-earphone', 'applied2', 'key']


def run_rakuten_search_at_shop_code(shop_code: str) -> None:
    logger.info('action=run_rakuten_search_at_shop_code status=run')

    rakuten = RakutenCrawler()
    rakuten.crawle_by_shop(shop_code)
    
    logger.info('action=run_rakuten_search_at_shop_code status=run')


def run_rakuten_search_all() -> None:
    logger.info('action=run_rakuten_search_all status=run')

    for shop_code in SHOP_CODES:
        try:
            run_rakuten_search_at_shop_code(shop_code)
        except Exception as ex:
            logger.error({'message': ex})

    logger.info('action=run_rakuten_search_all status=done')
