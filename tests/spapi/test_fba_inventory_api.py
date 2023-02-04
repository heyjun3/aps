import os
import json

from spapi.fba_inventory_api import FBAInventoryAPIParser


dirname = os.path.join(os.path.dirname(__file__), "test_json")


class TestFBAInventoryAPI(object):

    def test_parse_resp(self):
        path = os.path.join(dirname, "fba_inventory_api_v1.json")
        with open(path, "r") as f:
            file = f.read()

        products = FBAInventoryAPIParser.parse_fba_inventory_api_v1(json.loads(file))

        assert products[0] == {"sku": "4969133908026-N-2970-221129", "fnsku": "X00132SKM9"}
        assert products[-1] == {"sku": "4986773181114-N-7980-221210", "fnsku": "X0012O9JJB"}
