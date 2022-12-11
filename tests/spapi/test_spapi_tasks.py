from spapi.spapi_tasks import RunAmzTask
from spapi.models import AsinsInfo
from mws.models import MWS



class TestSpapiTasks(object):

    def test_validation_parameter(self):
        param = {
            "cost": 1111, "jan": "41922222",
            "filename": "filename", "url": "https://google.com"}
        client = RunAmzTask()
        result = client._validation_parameter(param)
        assert result == param

    def test_validation_parameter_not_in_cost(self):
        param = {
            "jan": "41922222",
            "filename": "filename", "url": "https://google.com"}
        client = RunAmzTask()
        result = client._validation_parameter(param)
        assert result == None

    def test_validation_parameter_not_in_jan(self):
        param = {
            "cost": 41922222,
            "filename": "filename", "url": "https://google.com"}
        client = RunAmzTask()
        result = client._validation_parameter(param)
        assert result == None

    def test_validation_parameter_not_in_filename(self):
        param = {
            "jan": "41922222",
            "cost": 1111, "url": "https://google.com"}
        client = RunAmzTask()
        result = client._validation_parameter(param)
        assert result == None

    def test_validation_parameter_not_in_url(self):
        param = {
            "jan": "41922222", "cost": 1111,
            "filename": "filename"}
        client = RunAmzTask()
        result = client._validation_parameter(param)
        assert result == None

    def test_map_asin_info_and_message(self):
        messages = [
            {"cost": 1111, "jan": "4444", "filename": "test_1", "url": "https://google.com"},
            {"cost": 2222, "jan": "5555", "filename": "test_2", "url": "https://yahoo.co.jp"},
            {"cost": 3333, "jan": "1111", "filename": "test_3", "url": "https://yahoo.co.jp"},
            ]
        asins = [
            AsinsInfo("test1", "4444", "test_title_1", 1),
            AsinsInfo("test2", "5555", "test_title_2", 2),
            AsinsInfo("test3", "6666", "test_title_3", 1),
        ]
        client = RunAmzTask()
        send_messages, result = client._map_asin_info_and_message(messages, asins)
        assert len(result) == 2
        assert result[0].asin == "test1"
        assert result[0].filename == "test_1"
        assert result[0].title == "test_title_1"
        assert result[0].jan == "4444"
        assert result[0].unit == 1
        assert result[0].cost == 1111
        assert result[0].url == "https://google.com"
        assert len(send_messages) == 1
        assert send_messages == [{"cost": 3333, "jan": "1111", "filename": "test_3", "url": "https://yahoo.co.jp"}]