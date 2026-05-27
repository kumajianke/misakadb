import time

from network import sock


class commandSend:
    def __init__(self, cli:sock.clientCore) -> None:
        self.cli : sock.clientCore= cli

    
    def init_command(self):
        start_time = time.time()
        print("开始初始化服务信息...")

        server_recv = self.cli.send_command("get-service-info")
        end_time = time.time()
        print(f"初始化服务信息耗时: {(end_time - start_time) * 1000.0:.2f} ms. [1send 1recv]",end="\n" *3)
        print(server_recv, end="\n" *3)

        print("初始化服务信息完成, 按Enter键继续")
        input("")
        
        return server_recv
