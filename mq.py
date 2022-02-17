import json
import time

import pika
from pika.exceptions import AMQPConnectionError

import log_settings


logger = log_settings.get_logger(__name__)


class MQ(object):

    def __init__(self, queue_name: str):
        self.queue_name = queue_name
        self.channel = self.create_mq_channel()
        self.properties = pika.BasicProperties(delivery_mode=pika.spec.PERSISTENT_DELIVERY_MODE)

    def create_mq_channel(self) -> pika.BaseConnection.channel:
        logger.info('action=create_mq_channel status=run')

        connection = pika.BlockingConnection()
        channel = connection.channel()
        channel.queue_declare(self.queue_name, durable=True)

        logger.info('action=create_mq_channel status=done')
        return channel

    def publish(self, value: str):
        logger.info('action=publish status=run')
        try:
            self.channel.basic_publish(exchange='', routing_key=self.queue_name, body=value, properties=self.properties)
        except AMQPConnectionError as ex:
            self.channel = self.create_mq_channel()
            self.channel.basic_publish(exchange='', routing_key=self.queue_name, body=value, properties=self.properties)

        logger.info('action=publish status=done')
    
    def receive(self):

        while True:
            resp = [self.channel.basic_get(self.queue_name, auto_ack=True) for _ in range(5) if self.channel.queue_declare(self.queue_name, durable=True).method.message_count]
            if resp:
                asin_list = list(map(lambda x: x[2].decode(), resp))
                yield asin_list
            else:
                time.sleep(60)


if __name__ == '__main__':
    mq = MQ('mws')
    a = mq.receive()
    for i in a:
        print(i)