import argparse
import sys
import getpass
import time

from interface.status import StatusSocket
from mql.MQ import MiQL
from network.command_send import commandSend
import network.sock as sock


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="MisakaDB CLI V 0.0.1")
    parser.add_argument("--address", type=str, default="127.0.0.1", help="服务地址, 默认本地")
    parser.add_argument("--port", type=int, default=10032, help="服务端口, 默认10032")
    parser.add_argument("--mode", type=str, default="default", help="运行模式: shell[cli调试模式]、onlyConn[仅连接]")
    parser.add_argument("--username", type=str, default="", help="[可选]用户名")
    parser.add_argument("--password", type=str, default="", help="[可选]密码")
    

    args = parser.parse_args()
    
    address = args.address
    port = args.port
    mode = args.mode
    username = args.username
    password = args.password
    
    
    cli = sock.clientCore(address, port)
    connect_err = cli.connect()
    if connect_err is not None:
        print(f"连接失败: {connect_err}", file=sys.stderr)
        cli.s.close()
        sys.exit(1)

    command_send = commandSend(cli)
    server_recv = command_send.init_command()

    if mode in ["shell", "sh", "s"]:
        retry = False
        while username == "" or password == "":
            if retry:
                print("\n登录失败, 请重新输入用户名和密码:\n")
            else:
                retry = True
            username = input("请输入用户名: ")
            password = getpass.getpass("请输入密码: ")
        
        
        
        while not command_send.login(username, password):
            print("登录失败, 请重新输入用户名和密码:\n")
            username = input("请输入用户名: ")
            password = getpass.getpass("请输入密码: ")
        
        cli.status = StatusSocket.Connected
        
        
        while True:
            commands = ""

            
            user_input = "\\"
            row = 0 
            while user_input.endswith("\\"):
                user_input = input(f"[{address}:{port} {cli.status.value}]{row if row > 0 else ''} >")
                commands += user_input + "\n"
                row += 1

            commands = commands.strip().replace("\\", " ")
            mq = MiQL(cli)
            if commands.startswith("mq."):
                try:
                    start_time = time.time()
                    res = mq.shot(commands)
                    end_time = time.time()
                    elapsed_ms = (end_time - start_time) * 1000.0

                    if isinstance(res, bytes):
                        res_str = res.decode("utf-8", errors="replace")
                        command_send.handle_response(res_str, elapsed_ms, is_init=False)
                    elif res is not None:
                        print(res)
                except Exception as e:
                    print(f"MiQL 执行失败: {e}", file=sys.stderr)
                continue

            if commands in ["exit", 'e', 'q']:
                cli.s.close()
                print("bye~")
                sys.exit(0)

            if commands:
                start_time = time.time()
                res = cli.send_command(commands)
                end_time = time.time()
                elapsed_ms = (end_time - start_time) * 1000.0
                
                res_str = res.decode('utf-8', errors='replace')
                command_send.handle_response(res_str, elapsed_ms, is_init=False)

        
