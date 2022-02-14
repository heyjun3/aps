
import pika


connection = pika.BlockingConnection()
channel = connection.channel()
channel.queue_declare(queue='hello')


def send():
    channel.basic_publish(exchange='', routing_key='hello', body='hello World')
    print('send')


def receive():

    def callback(ch, method, properties, body):
        print(f'Received {body}')

    channel.basic_consume(queue='hello', on_message_callback=callback, auto_ack=True)
    channel.start_consuming()


if __name__ == '__main__':
    send()