import json
import time
from types import FunctionType
import threading

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
    
    def receive(self) -> None:

        while True:
            resp = [self.channel.basic_get(self.queue_name, auto_ack=True) for _ in range(5) if self.channel.queue_declare(self.queue_name, durable=True).method.message_count]
            if resp:
                asin_list = list(map(lambda x: x[2].decode(), resp))
                yield asin_list
            else:
                yield None
                time.sleep(30)

    def get(self, interval_sec: int=3) -> str:

        while True:
            resp = self.channel.basic_get(self.queue_name, auto_ack=True)
            _method, _properties, body = resp
            if body:
                yield body.decode()
            else:
                time.sleep(interval_sec)

    def callback_recieve(self, func: FunctionType, interval_sec: float=1.0) -> None:
        logger.info('action=run_callback_recieve status=run')

        def callback(ch, method, properties, body):
            thread = threading.Thread(target=func, args=(json.loads(body.decode()),))
            thread.start()
            time.sleep(interval_sec)

        self.channel.basic_consume(queue=self.queue_name, on_message_callback=callback, auto_ack=True)
        self.channel.start_consuming()

        logger.info('action=run_callback_recieve status=done')

def emit_log() -> None:
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host='localhost'))
    
    channel = connection.channel()
    channel.exchange_declare(exchange='logs', exchange_type='fanout')

    message = 'info: hello world!'
    channel.basic_publish(exchange='logs', routing_key='', body=message)
    connection.close()


def receive_logs(name: str) -> None:
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host="localhost"))

    channel = connection.channel()
    channel.exchange_declare(exchange="logs", exchange_type='fanout')
    result = channel.queue_declare(queue=name, exclusive=True)
    queue_name = result.method.queue

    channel.queue_bind(exchange="logs", queue=queue_name)
    while True:
        resp = channel.basic_get(queue_name, auto_ack=True)
        if resp[2]:
            print(resp)
        else:
            time.sleep(10)

if __name__ == '__main__':
    # receive_logs('two')
    emit_log()