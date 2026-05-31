from enum import Enum


class StatusSocket(Enum):
    Disconnected = "Disconnected"
    ConnectedNoAuth = "Connected(NoAuth)"
    Connected = "Connected(Auth)"
   