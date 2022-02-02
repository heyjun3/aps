from mws.api import AmazonClient


mws = AmazonClient()
p = 'B09L1HJ15J'
a = mws.get_lowest_priced_offers_for_asin(p)
