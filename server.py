from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

from spapi.fba_inventory_api import FBAInventoryAPI
from spapi.listings_items_api import ListingsItemsAPI
from spapi.spapi import TooMatchParameterException, BadItemTypeException

app = FastAPI()
client = FBAInventoryAPI()
listing_client = ListingsItemsAPI()


@app.get("/inventory-summaries")
async def get_inventory_summaries(next_token: str = ''):
    return await client.get_inventory_summaries(next_token)


@app.get('/get-pricing')
async def get_pricing(ids: str, id_type: str):
    try:
        return await client.get_pricing(ids.split(','), id_type)
    except TooMatchParameterException:
        raise HTTPException(
            status_code=400, detail="too many ids. expect less than 20 ids")
    except BadItemTypeException:
        raise HTTPException(
            status_code=400, detail="Bad item type. expect Asin or Sku")
    except Exception:
        raise HTTPException(status_code=503, detail="Internal Server Error")
    

@app.get('/get-listing-offers-batch')
async def get_listing_offers_batch(skus: str):
    return await client.get_listing_offers_batch(skus.split(','))

class UpdatePriceInput(BaseModel):
    sku: str
    price: int

@app.post('/price')
async def update_price(input: UpdatePriceInput):
    return await listing_client.update_price(input.sku, input.price)
