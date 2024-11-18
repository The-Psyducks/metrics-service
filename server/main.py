import pika, os, signal, sys
from http.server import BaseHTTPRequestHandler, HTTPServer

# Access the CLODUAMQP_URL environment variable and parse it (fallback to localhost)
print("Connecting to CloudAMQP...")
url = os.environ.get('CLOUDAMQP_URL', 'amqp://admin:admin123@rabbitmq:5672/%2f')
params = pika.URLParameters(url)
connection = pika.BlockingConnection(params)
print("creating channel...")
channel = connection.channel() # start a channel
channel.queue_declare(queue='hello') # Declare a queue
print("publishing...")
channel.basic_publish(exchange='',
                  routing_key='hello',
                  body='Hello CloudAMQP!')

print(" [x] Sent 'Hello World!'")

def callback(ch, method, properties, body):
  print(" [x] Received " + str(body))

print(' [*] Waiting for messages:')
channel.basic_consume('hello',
                      callback,
                      auto_ack=True)

print(' [*] Waiting for messages:')
# channel.start_consuming()
# channel.consume('hello', callback, auto_ack=True)
connection.close()

print("Holaaaa")

# HTTP server
class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()
        self.wfile.write(b"Twitsnap RabbitMQ Consumer is running!")

# Bind to the port specified in the environment variable PORT, or default to 8080
port = int(os.environ.get('PORT', 8080))
server_address = ('', port)

print(f"Starting HTTP server on port {port}...")
httpd = HTTPServer(server_address, RequestHandler)

try:
    print(' [*] Waiting for messages and serving HTTP:')
    httpd.serve_forever()
except KeyboardInterrupt:
    print(" Shutting down HTTP server...")
    httpd.server_close()
finally:
    channel.stop_consuming()
    connection.close()