
import time

from apis.api import MisakaDBClient


clients = MisakaDBClient("127.0.0.1", 10032)
clients.connect()

clients2 = MisakaDBClient("127.0.0.1", 10032)
clients2.connect()


clients3 = MisakaDBClient("127.0.0.1", 10032)
clients3.connect()


time.sleep(3)
print(clients.heart.stats())
print(clients2.heart.stats())
print(clients3.heart.stats())
    