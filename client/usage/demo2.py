import time

from apis.api import MisakaDBClient


clients = MisakaDBClient("127.0.0.1", 10032)
clients.connect()
time.sleep(10)
print(clients.heart.stats())
