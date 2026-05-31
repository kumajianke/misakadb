import json
import time
from typing import Any

from network import sock

MISAKA_LOGO = [
    "             .M.             .M.             ",
    "            .MMM.           .MMM.            ",
    "           .MMMMM.         .MMMMM.           ",
    "          .MMMMMMM.       .MMMMMMM.          ",
    "         .MMMMMMMMM.     .MMMMMMMMM.         ",
    "        .MMMMMMMMMMM.   .MMMMMMMMMMM.        ",
    "       .MMMMMMMMMMMMM. .MMMMMMMMMMMMM.       ",
    "      .MMMMMMMMMMMMMMM.MMMMMMMMMMMMMMM.      ",
    "     .MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM.     ",
    "    .MMMMMMMMM...MMMMMMMMMMM...MMMMMMMMM.    ",
    "   .MMMMMMMMM.   .MMMMMMMMM.   .MMMMMMMMM.   ",
    "  .MMMMMMMMM.     .MMMMMMM.     .MMMMMMMMM.  ",
    " .MMMMMMMMM.       .MMMMM.       .MMMMMMMMM. ",
    ".MMMMMMMMM.         .MMM.         .MMMMMMMMM.",
]

DEFAULT_LABELS = {
    "network.address": "监听地址",
    "network.port": "监听端口",
    "network.max_conn": "最大连接数",
    "network.retry_count": "重试次数",
    "network.retry_delay": "重试延迟(秒)",
    "service.version": "服务版本",
}

DEFAULT_ORDER = [
    "service.version",
    "network.address",
    "network.port",
    "network.max_conn",
    "network.retry_count",
    "network.retry_delay",
]


class commandSend:
    def __init__(self, cli:sock.clientCore) -> None:
        self.cli : sock.clientCore= cli

    def _flatten_json(self, data: dict[str, Any], prefix: str = "") -> list[tuple[str, Any]]:
        items: list[tuple[str, Any]] = []
        for key, value in data.items():
            if key == "_display":
                continue

            full_key = f"{prefix}.{key}" if prefix else key
            if isinstance(value, dict):
                items.extend(self._flatten_json(value, full_key))
            else:
                items.append((full_key, value))
        return items

    def _format_value(self, value: Any) -> str:
        if isinstance(value, list):
            return ", ".join(str(item) for item in value)
        if isinstance(value, (dict, tuple)):
            return json.dumps(value, ensure_ascii=False)
        return str(value)

    def _build_info_lines(self, json_data: dict[str, Any], elapsed_ms: float) -> list[str]:
        display_config = json_data.get("_display", {})
        custom_labels = display_config.get("labels", {})
        custom_order = display_config.get("order", [])

        labels = {**DEFAULT_LABELS, **custom_labels}
        ordered_keys = custom_order or DEFAULT_ORDER
        flattened_items = self._flatten_json(json_data)
        flattened_map = dict(flattened_items)

        ordered_items: list[tuple[str, Any]] = []
        used_keys: set[str] = set()

        for key in ordered_keys:
            if key in flattened_map:
                ordered_items.append((key, flattened_map[key]))
                used_keys.add(key)

        for key, value in flattened_items:
            if key not in used_keys:
                ordered_items.append((key, value))

        def get_display_width(text):
            return sum(2 if ord(c) > 127 else 1 for c in text)

        info_lines = [
            "\033[1;36m  库码科技工作室",
            "\033[1;32m  MisakaDB Client\033[0m",
            "",
            f"\033[1;36m请求耗时\033[{'1;32m' if round(elapsed_ms, 2) < 0.2 else '1;33m'}{' ' * 6} : {elapsed_ms:.2f} ms",
        ]

        for key, value in ordered_items:
            label = labels.get(key, key)
            label_width = get_display_width(label)
            padding = max(0, 14 - label_width)
            padded_label = label + " " * padding
            info_lines.append(f"\033[1;36m{padded_label}\033[0m : {self._format_value(value)}")

        return info_lines

    def _print_fetch_layout(self, info_lines: list[str]) -> None:
        logo_width = max(len(line) for line in MISAKA_LOGO) if MISAKA_LOGO else 0
        total_lines = max(len(info_lines), len(MISAKA_LOGO))

        for index in range(total_lines):
            logo_part = MISAKA_LOGO[index] if index < len(MISAKA_LOGO) else ""
            info_part = info_lines[index] if index < len(info_lines) else ""
            
            if logo_part:
                left_aligned = f"\033[1;36m{logo_part:<{logo_width}}\033[0m"
            else:
                left_aligned = " " * logo_width
                
            print(f"{left_aligned}    {info_part}")

    
    def handle_response(self, server_recv_str: str, elapsed_ms: float, is_init: bool = False) -> str:
        if is_init:
            if server_recv_str.startswith("[err]"):
                info_lines = [
                    "\033[1;36m  库码科技工作室",
                    "\033[1;31m  MisakaDB Client (Failed)\033[0m",
                    "",
                    f"\033[1;36m请求耗时\033[1;33m{' ' * 6} : {elapsed_ms:.2f} ms",
                    f"\033[1;31m错误信息\033[0m : {server_recv_str[5:]}"
                ]
                self._print_fetch_layout(info_lines)
                print()
                return server_recv_str

            elif server_recv_str.startswith("[ok]"):
                try:
                    json_data = json.loads(server_recv_str[4:])
                    self._print_fetch_layout(self._build_info_lines(json_data, elapsed_ms))
                except Exception:
                    # 兼容非JSON格式的普通响应
                    info_lines = [
                        "\033[1;36m  库码科技工作室",
                        "\033[1;32m  MisakaDB Client\033[0m",
                        "",
                        f"\033[1;36m请求耗时\033[{'1;32m' if round(elapsed_ms, 2) < 0.2 else '1;33m'}{' ' * 6} : {elapsed_ms:.2f} ms",
                        f"\033[1;36m响应内容\033[0m : {server_recv_str[4:]}"
                    ]
                    self._print_fetch_layout(info_lines)
                print()
                return server_recv_str
            else:
                print(server_recv_str)
                return server_recv_str
        else:
            # 普通命令响应格式
            time_color = '\033[1;32m' if round(elapsed_ms, 2) < 0.2 else '\033[1;33m'
            if server_recv_str.startswith("[err]"):
                print(f"\033[1;31mError: {server_recv_str[5:]}\033[0m ({time_color}{elapsed_ms:.2f} ms\033[0m)")
            elif server_recv_str.startswith("[error]"):
                print(f"\033[1;31mError: {server_recv_str[7:]}\033[0m ({time_color}{elapsed_ms:.2f} ms\033[0m)")
            elif server_recv_str.startswith("[ok]"):
                print(f"{server_recv_str[4:]} ({time_color}{elapsed_ms:.2f} ms\033[0m)")
            else:
                print(f"{server_recv_str} ({time_color}{elapsed_ms:.2f} ms\033[0m)")
            return server_recv_str

    def init_command(self):
        start_time = time.time()
        print("开始初始化服务信息...")

        server_recv_bytes = self.cli.send_command("get-service-info")
        end_time = time.time()
        elapsed_ms = (end_time - start_time) * 1000.0

        try:
            server_recv_str = server_recv_bytes.decode('utf-8')
        except Exception:
            server_recv_str = str(server_recv_bytes)

        self.handle_response(server_recv_str, elapsed_ms, is_init=True)

        print("初始化服务信息完成, 按Enter键继续")
        input("")
            
        return server_recv_str

    
    def login(self, username: str, password: str):
        res = b""
        start_time = time.time()
        if username and password:
            res = self.cli.send_command(f"login {username} {password}")
            
        # 确保res可以被decode，如果为空则设为空字符串的bytes
        if not res:
            res = b""
            
        res_str = res.decode("utf-8", errors="replace")
        end_time = time.time()
        elapsed_ms = (end_time - start_time) * 1000.0
        
        if res_str:
            self.handle_response(res_str, elapsed_ms)
            
        return res_str.startswith("[ok]")
