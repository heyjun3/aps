import json
import pika
import time


connection = pika.BlockingConnection()
channel = connection.channel()
channel.queue_declare(queue='hello')


def send():
    channel.basic_publish(exchange='', routing_key='hello', body=json.dumps({'hello': 1}))
    print('send')


def receive():

    # print(type(channel.queue_declare('hell').method.message_count))
    # print(q.method.message_count)

    # for method_frame, properties, body in channel.consume('hello', auto_ack=True):
    #     print(body)
    # while True:
    #     value, a, b = channel.consume('hello', auto_ack=True)
    #     print(value)
    # value = channel.basic_get('hello', auto_ack=True)
    # print(value)
    # channel.basic_consume('hello', on_message_callback=lambda x: print(x), auto_ack=True)
    # channel.start_consuming()
    while True:
        resp = [channel.basic_get('hello', auto_ack=True) for _ in range(5) if channel.queue_declare('hello').method.message_count]
        asin_list = list(map(lambda x: x[2].decode(), resp))
        print(asin_list)
        time.sleep(60)
    # print(type(json.loads(asin_list[0]).get('hello')))

if __name__ == '__main__':
    q = channel.queue_declare(queue='hello')
    print(type(q))
    # print(q.method.message_count)
    # receive()