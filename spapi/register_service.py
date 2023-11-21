import time
import datetime
from pathlib import Path
from typing import List
import csv
import io
import gzip

import requests
import gspread
from oauth2client.service_account import ServiceAccountCredentials

import settings
import log_settings
from spapi.listings_items_api import ListingsItemsAPI
from spapi.fba_inventory_api import FBAInventoryAPI
from spapi.fba_inventory_api import FBAInventoryAPIParser
from spapi.feeds_api import FeedsAPI


logger = log_settings.get_logger(__name__)
SCOPE = ('https://spreadsheets.google.com/feeds',
         'https://www.googleapis.com/auth/drive')


class RegisterService(object):

    def __init__(self, credential_file: str) -> None:
        self.client = self._create_sheet_client(credential_file)
        self.spapi = ListingsItemsAPI()
        self.inventory = FBAInventoryAPI()
        self.feed_client = FeedsAPI()

    # TODO レコードのバリデーション追加する。
    async def start_register(self, title: str, name: str, interval_sec: int = 2):
        logger.info({"action": "start_register", "status": "run"})
        add = self.client.open(title).worksheet(name)
        records = add.get_all_records()

        self._validate_records(records)
        for record in records:
            sku = self._generate_sku(record)
            res = await self.spapi.create_new_sku(sku, record.get("PRICE"),
                                                  record.get("ASIN"), settings.CONDITION_NOTE)
            time.sleep(interval_sec)
            if res.get("status") == "ACCEPTED":
                logger.info({"action": "start_register", "message": "register request is accepted",
                             "sku": sku})
                continue
            logger.error({"action": "start_register", "message": "register request is failed",
                          "value": record, "response": res})

        logger.info({"action": "start_register", "status": "done"})

    async def check_registerd(self, title: str, get_sheet: str, keep_sheet: str):
        logger.info({"action": "check_registerd", "status": "run"})

        add = self.client.open(title).worksheet(get_sheet)
        records = add.get_all_records()

        self._validate_records(records)
        for record in records:
            record["SKU"] = self._generate_sku(record)
        inventories = []
        for i in range(0, len(records), 50):
            skus = [record.get("SKU") for record in records[i:i+50]]
            res = await self.inventory.fba_inventory_api_v1(skus)
            inventory = FBAInventoryAPIParser.parse_fba_inventory_api_v1(res)
            time.sleep(10)
            if inventory:
                inventories.extend(inventory)

        fnsku = {inv["sku"]: inv["fnsku"] for inv in inventories}
        for record in records:
            record["FNSKU"] = fnsku.get(record.get("SKU"))

        rows = [list(record.values())
                for record in records if record.get("FNSKU") is not None]
        db = self.client.open(title).worksheet(keep_sheet)
        db.append_rows(rows)

        for record in records:
            fnsku = record.get("FNSKU")
            if not fnsku:
                continue
            cell = add.find(record.get("ASIN"), in_column=3)
            time.sleep(1)
            if cell:
                add.delete_row(cell.row)

        point_record = [[record.get('SKU'), record.get('POINT')] for record in records if record.get(
            'FNSKU') is not None and record.get('POINT') is not None]

        if point_record:
            await self.register_points(point_record)

        logger.info({"action": "check_registerd", "status": "done"})

    # INFO show spapi feeds usecase page
    # items = [[sku: str, point: int]]
    async def register_points(self, items: list[list], interval_sec: int = 1):
        logger.info({'action': '_register_points', 'status': 'run'})
        header = ['sku', 'points_percent']
        rows = [header, *list(filter(lambda x: int(x[1]) <= 100, items))]
        feed = io.StringIO()
        csv.writer(feed, delimiter='\t').writerows(rows)
        send_tsv = feed.getvalue().encode('UTF-8')

        res = await self.feed_client.create_feed_document('text/tsv', 'UTF-8')

        logger.info({'action': '_register_points', 'send_tsv': send_tsv})
        requests.put(res['url'], data=send_tsv, headers={
                     'Content-Type': 'text/tsv; charset=UTF-8'})

        r = await self.feed_client.create_feed('POST_FLAT_FILE_OFFER_POINTS_PREFERENCE_DATA', res['feedDocumentId'])
        while True:
            r = await self.feed_client.get_feed(r['feedId'])
            if r.get('processingStatus') == 'DONE':
                break
            time.sleep(interval_sec)

        res = await self.feed_client.get_feed_document(r['resultFeedDocumentId'])

        r = requests.get(res['url'], stream=True)
        gzip_file = io.BytesIO(r.content)
        with gzip.open(gzip_file, 'rt') as f:
            data = f.read()

        logger.info({'action': '_register_points',
                    'status': 'done', 'result': data})

    def _create_sheet_client(self, credential: str) -> gspread.Client:
        path = Path.cwd().joinpath(credential)
        keyfile = ServiceAccountCredentials.from_json_keyfile_name(path, SCOPE)
        return gspread.authorize(keyfile)

    def _validate_records(self, records: List[dict]) -> bool:
        for record in records:
            if not all([record.get(key) for key in ["NAME", "ASIN", "JAN", "DIVISION", "PRICE", "COST"]]):
                raise Exception("validation error")

        return True

    def _generate_sku(self, record: dict) -> str:
        date = datetime.datetime.now().strftime("%Y%m%d")
        return '-'.join([str(record.get(key)) for key in ["JAN", "DIVISION", "COST"]] + [date])
