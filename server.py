from fastapi import FastAPI

from spapi.fba_inventory_api import FBAInventoryAPI

app = FastAPI()
client = FBAInventoryAPI()

@app.get("/inventory-summaries")
async def get_inventory_summaries(next_token: str = ''):
    return await client.get_inventory_summaries(next_token)
