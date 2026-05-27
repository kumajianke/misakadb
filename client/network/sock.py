import socket
import time

from interface.status import StatusSocket

class clientCore:
    def __init__(self, server_address, server_port) -> None:
        self.server_address = server_address
        self.server_port = server_port
        # 创建 TCP 套接字
        self.s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.s.settimeout(5)
        self.status = StatusSocket.Disconnected

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
        self.s.send(data.encode("utf-8"))
    
    def recv_str(self):
        data = self.s.recv(1024)
        return data.decode("utf-8")
    
    def send_bytes(self, data:bytes):
        self.s.send(data)
        
    def recv_bytes(self):
        data = self.s.recv(1024)
        return data

    def send_command(self, command:str):
        self.send_str(command)
        return self.recv_bytes()
    
    
