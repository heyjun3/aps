from spapi.spapi_tasks import RunAmzTask



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
