import argparse
import asyncio
import sys

from network.command_send import commandSend
import network.sock as sock


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="MisakaDB CLI V 0.0.1")
    parser.add_argument("--address", type=str, default="127.0.0.1", help="服务地址, 默认本地")
    parser.add_argument("--port", type=int, default=10032, help="服务端口, 默认10032")
    parser.add_argument("--mode", type=str, default="default", help="运行模式: shell[cli调试模式]、onlyConn[仅连接]")

    args = parser.parse_args()
    
    address = args.address
    port = args.port
    mode = args.mode
    
    
    cli = sock.clientCore(address, port)
    connect_err = cli.connect()
    if connect_err is not None:
        print(f"连接失败: {connect_err}", file=sys.stderr)
        cli.s.close()
        sys.exit(1)

    command_send = commandSend(cli)
    server_recv = asyncio.run(command_send.init_command())

    if mode in ["shell", "sh", "s"]:
        while True:
            commands = ""

            
            user_input = "\\"
            row = 0 
            while user_input.endswith("\\"):
                user_input = input(f"[{address}:{port} {cli.status.value}]{row if row > 0 else ''} >")
                commands += user_input + "\n"
                row += 1

            commands = commands.strip().replace("\\", " ")

            if commands in ["exit", 'e', 'q']:
                cli.s.close()
                print("bye~")
                sys.exit(0)
        
