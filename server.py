from typing import List

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

from spapi.fba_inventory_api import FBAInventoryAPI
from spapi.listings_items_api import ListingsItemsAPI
from spapi.register_service import RegisterService
from spapi.spapi import TooMatchParameterException, BadItemTypeException
import settings

app = FastAPI()
client = FBAInventoryAPI()
listing_client = ListingsItemsAPI()
register_service = RegisterService(settings.CREDENTIAL_FILE_NAME)


@app.get("/inventory-summaries")
async def get_inventory_summaries(next_token: str = '', is_detail: bool = True):
    return await client.get_inventory_summaries(next_token, is_detail)


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

@app.get('/competitive-pricing')    
async def get_competitive_pricing(asins: str):
    return await client.get_competitive_pricing(asins.split(','))

@app.get('/get-listing-offers-batch')
async def get_listing_offers_batch(skus: str):
    return await client.get_listing_offers_batch(skus.split(','))

class UpdatePriceInput(BaseModel):
    sku: str
    price: int

@app.post('/price')
async def update_price(input: UpdatePriceInput):
    return await listing_client.update_price(input.sku, input.price)

class UpdatePointInput(BaseModel):
    sku: str
    percent_point: int


@app.post('/points')
async def update_point(input: List[UpdatePointInput]):
    items = [[item.sku, item.percent_point] for item in input]
    return await register_service.register_points(items)
