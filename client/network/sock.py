import random
import socket
import threading
import time
from uuid import uuid4

from interface.status import StatusSocket

class clientCore:
    def __init__(self, server_address, server_port) -> None:
        self.server_address = server_address
        self.server_port = server_port
        # 创建 TCP 套接字
        self.s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.s.settimeout(5)
        self.status = StatusSocket.Disconnected
        
        self.lock = threading.Lock()
        
        
        

        self.mode = []
        
    def connect(self):
        last_err = None
        # 重试 3 次
        for i in range(3):
            try:
                self.s.connect((self.server_address, self.server_port))
                self.status = StatusSocket.Connected
                return None # 连接成功返回 None
            except socket.error as e:
                last_err = e
                print(f"连接失败，第 {i + 1} 次重试...")
                if i < 2:
                    time.sleep(1) # 失败后等待 1 秒再重试
                    
        return last_err # 3 次都失败，返回最后的错误

    def send_str(self, data:str):
        
        data_bytes = data.encode("utf-8")
        self.send_bytes(data_bytes)
    
    def recv_str(self):
        data = self.recv_bytes()
        return data.decode("utf-8")
    
    def send_bytes(self, data:bytes, msg_type:int=0x00):
        seq_id = random.randint(10000, 100000)
        seq_id_bytes = seq_id.to_bytes(4, byteorder='big')
        
        # 4-byte length (Big Endian)
        length_bytes = len(data).to_bytes(4, byteorder='big')
        # 1-byte message type
        type_byte = msg_type.to_bytes(1, byteorder='big')
        # send length + type + data
        with self.lock:
            
            self.s.sendall(seq_id_bytes + length_bytes + type_byte + data)
        
    def send_heartbeat(self):
        # 发送心跳包: 数据长度为0，类型为0x01

        self.send_bytes(b"", msg_type=0x01)
        
    def recv_bytes(self):
        # read 4-byte length
        length_bytes = self._recv_exact(4)
        if not length_bytes:
            return b""
        length = int.from_bytes(length_bytes, byteorder='big')
        
        # read length bytes of data
        if length > 0:
            return self._recv_exact(length)
        return b""

    def _recv_exact(self, n: int) -> bytes:
        data = bytearray()
        while len(data) < n:
            packet = self.s.recv(n - len(data))
            if not packet:
                return b""
            data.extend(packet)
        return bytes(data)

    def send_command(self, command:str):
        self.send_str(command)
        return self.recv_bytes()
    
    
