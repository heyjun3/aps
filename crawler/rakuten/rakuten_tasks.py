from cmath import log
import log_settings
from crawler.rakuten.rakuten import RakutenAPIClient, RakutenCrawler

logger = log_settings.get_logger(__name__)


SHOP_CODES = [
    'superdeal', 'ksdenki', 'dj', 'e-zoa', 'reckb', 'jtus', 'ioplaza', 'ikebe',
    'shoptsukumo', 'premiumgt', 'acer-direct', 'wakeari', 'jism', 'sakurayama',
    'ikeshibu', 'aikyoku-bargain-center', 'aikyoku', 'ikebe-rockhouse', 'ishibashi',
    'pckoubou', 'e-earphone', 'applied2', 'yamada-denki', 'key', 'r-kojima']


def run_rakuten_search_at_shop_code(shop_code: str) -> None:
    logger.info('action=run_rakuten_search_at_shop_code status=run')

    rakuten = RakutenCrawler(shop_code=shop_code)
    rakuten.main()
    
    logger.info('action=run_rakuten_search_at_shop_code status=run')


def run_rakuten_search_all() -> None:
    logger.info('action=run_rakuten_search_all status=run')

    for shop_code in SHOP_CODES:
        # run_rakuten_search_at_shop_code(shop_code)
        RakutenCrawler(shop_code).main()

    logger.info('action=run_rakuten_search_all status=done')
