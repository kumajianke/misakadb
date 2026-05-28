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
            "\033[1;36m Kumare @\033[1;36m MisakaDB Client\033[0m",
            "----------------",
            f"\033[1;36m请求耗时\033[0m{' ' * 6} : {elapsed_ms:.2f} ms",
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

    
    def init_command(self):
        start_time = time.time()
        print("开始初始化服务信息...")

        server_recv = self.cli.send_command("get-service-info")
        end_time = time.time()
        elapsed_ms = (end_time - start_time) * 1000.0

        if (server_recv.decode().startswith("error")) :
            print(f"初始化服务信息失败: {server_recv}")
            return server_recv

        json_data = json.loads(server_recv)
        self._print_fetch_layout(self._build_info_lines(json_data, elapsed_ms))
        print()

        print("初始化服务信息完成, 按Enter键继续")
        input("")
        
        return server_recv
