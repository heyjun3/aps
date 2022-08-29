import time


class Cache(object):

    def __init__(self, value: object, ttl_sec: int) -> None:
        self.ttl_sec = ttl_sec
        self.value = value
        self.start_time = time.time()

    def set_value(self, value: object) -> None:
        self.start_time = time.time()
        self.value = value

    def get_value(self) -> object:
        if time.time() > (self.start_time + self.ttl_sec):
            self.value = None
        return self.value
