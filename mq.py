import json
import time
from types import FunctionType
from typing import List
import threading

import pika
from pika.exceptions import AMQPConnectionError
from pika.exceptions import DuplicateGetOkCallback

import log_settings
import settings


logger = log_settings.get_logger(__name__)


class MQ(object):

    def __init__(self, queue_name: str):
        self.queue_name = queue_name
        self.queue = None
        self.credentials = pika.PlainCredentials(settings.MQ_USER, settings.MQ_PASSWORD)
        self.channel = self.create_mq_channel()
        self.properties = pika.BasicProperties(delivery_mode=pika.spec.PERSISTENT_DELIVERY_MODE)

    def create_mq_channel(self, retry_count: int=30, interval_sec=10) -> pika.BaseConnection.channel:
        logger.info('action=create_mq_channel status=run')

        for _ in range(retry_count):
            try:
                connection = pika.BlockingConnection(
                    pika.ConnectionParameters(
                        host=settings.MQ_HOST,
                        port=settings.MQ_PORT,
                        credentials=self.credentials,
                    )
                )
            except AMQPConnectionError as ex:
                logger.info({'error': ex})
                time.sleep(interval_sec)
                continue

            channel = connection.channel()
            self.queue = channel.queue_declare(self.queue_name, durable=True)

            logger.info('action=create_mq_channel status=done')
            return channel
        logger.error({'message': 'create mq channel failed. connection refused'})

    def get_message_count(self) -> int:
        logger.info('action=get_message_count status=run')
        while True:
            try:
                queue = self.channel.queue_declare(self.queue_name, durable=True)
                message_count = queue.method.message_count
            except AMQPConnectionError as ex:
                logger.error({'message': ex})
                self.create_mq_channel()
                continue
            
            logger.info({"action": "get_message_count", "count": message_count})
            return message_count

    def publish(self, value: str):
        logger.debug('action=publish status=run')
        try:
            self.channel.basic_publish(exchange='', routing_key=self.queue_name, body=value, properties=self.properties)
        except AMQPConnectionError as ex:
            logger.error({'message': ex})
            self.channel = self.create_mq_channel()
            self.channel.basic_publish(exchange='', routing_key=self.queue_name, body=value, properties=self.properties)
        except FileNotFoundError as ex:
            logger.error({'message': ex})

        logger.debug('action=publish status=done')
    
    def receive(self, get_count: int=20) -> List[str]:

        while True:
            messages = []
            try:
                for _ in range(get_count):
                    message = self.channel.basic_get(self.queue_name, auto_ack=True)
                    if not all(message):
                        break
                    _, _, body = message
                    messages.append(body.decode())
                # messages = [self.channel.basic_get(self.queue_name, auto_ack=True) for _ in range(get_count)]
            except (AMQPConnectionError, FileNotFoundError) as ex:
                logger.error({'message': ex})
                self.channel = self.create_mq_channel()

            if messages:
                yield messages
            else:
                yield None

    def get(self) -> str|None:

        while True:
            try:
                resp = self.channel.basic_get(self.queue_name, auto_ack=True)
            except (FileNotFoundError, DuplicateGetOkCallback) as e:
                logger.error({'message': e})
                yield None
                time.sleep(1)
                continue

            _method, _properties, body = resp
            if body:
                yield body.decode()
            else:
                yield None

    def basic_get(self) -> str|None:
        try:
            message = self.channel.basic_get(self.queue_name, auto_ack=True)
        except AMQPConnectionError as ex:
            logger.error({'mesage': ex})
            return None
        _method, _properties, body = message
        if body:
            return body.decode()
        else:
            return None

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
    channel.queue_declare('one', durable=True)
    channel.queue_declare('two', durable=True)
    channel.queue_bind(exchange="logs", queue="one")
    channel.queue_bind(exchange="logs", queue="two")

    message = 'info: hello world!'
    channel.basic_publish(exchange='logs', routing_key='', body=message)
    connection.close()


def receive_logs(queue_name: str) -> None:
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host="localhost"))

    channel = connection.channel()
    while True:
        resp = channel.basic_get(queue_name, auto_ack=True)
        if resp[2]:
            print(resp)
        else:
            time.sleep(10)

if __name__ == '__main__':
    mq = MQ('mws')
    print(type(mq.get_message_count()))