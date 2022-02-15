import json
import pika


connection = pika.BlockingConnection()
channel = connection.channel()
channel.queue_declare(queue='hello')


def send():
    channel.basic_publish(exchange='', routing_key='hello', body=json.dumps({'hello': 1}))
    print('send')


def receive():

    # for method_frame, properties, body in channel.consume('hello', auto_ack=True):
    #     print(body)
    # while True:
    #     value = channel.consume('hello', auto_ack=True)
    #     print(value)
    # value = [channel.basic_get('hello', auto_ack=True) for _ in range(5) if channel.basic_get('hello', auto_ack=True)[2] is not None]
    method_frame, header, body = channel.basic_get('hello', auto_ack=True)
    print(body)


if __name__ == '__main__':
    receive()
    # send()