# MisakaDB Python Client

这是 MisakaDB 的 Python 客户端文档，主要面向客户端调用方。

当前核心入口位于 `apis/api.py`，主要提供 `MisakaDBClient` 和心跳控制对象 `client.heart`。

## 功能概览

- 连接 MisakaDB 服务端
- 获取服务信息
- 执行文本命令
- 查询服务版本、网络配置、允许命令
- 自动发送心跳包
- 手动启动、停止心跳
- 查询心跳次数、失败次数、丢包率

## 目录

- 客户端 API: `apis/api.py`
- 命令行入口: `main.py`
- 底层 socket 实现: `network/sock.py`
- 压测示例: `usage/demo1.py`

## 快速开始

### 1. 创建客户端

```python
from apis.api import MisakaDBClient

client = MisakaDBClient("127.0.0.1", 10032)
```

默认参数：

- `host="127.0.0.1"`
- `port=10032`
- `heartbeat_interval=10.0`

### 2. 连接服务端

```python
ok = client.connect()
if not ok:
    print("连接失败")
```

也可以指定重试参数：

```python
ok = client.connect(retries=3, retry_delay=1.0)
```

### 3. 获取服务信息

```python
service_info = client.get_service_info()
print(service_info)
```

典型返回值：

```python
{
    "network": {
        "port": 10032,
        "address": "0.0.0.0",
        "max_conn": 2,
        "retry_count": 3,
        "retry_delay": 30
    },
    "service": {
        "version": "0.0.3",
        "allow_command": ["get-service-info"]
    }
}
```

### 4. 关闭连接

```python
client.close()
```

## 推荐写法

推荐使用上下文管理器，避免忘记关闭连接：

```python
from apis.api import MisakaDBClient

with MisakaDBClient("127.0.0.1", 10032) as client:
    info = client.get_service_info()
    print(info)
```

## API 说明

### `MisakaDBClient(host="127.0.0.1", port=10032, heartbeat_interval=10.0)`

创建客户端对象。

参数说明：

- `host`: 服务端地址
- `port`: 服务端端口
- `heartbeat_interval`: 心跳发送间隔，单位秒

注意：

- 初始化时会创建 `client.heart` 心跳控制器
- 初始化时会启动心跳线程
- 真正发送心跳前仍要求连接已建立

### `connect(retries=3, retry_delay=1.0) -> bool`

连接到服务端。

参数说明：

- `retries`: 失败重试次数
- `retry_delay`: 重试间隔，单位秒

返回值：

- `True`: 连接成功
- `False`: 连接失败

### `close() -> None`

关闭当前连接，并停止当前心跳控制器。

### `get_service_info() -> dict | None`

获取服务端信息。

成功时返回字典，失败时返回 `None`。

### `execute_command(command: str) -> str | dict | list | None`

发送一条文本命令到服务端。

行为说明：

- 如果返回内容是 JSON，会自动解析成 `dict` 或 `list`
- 如果返回内容不是 JSON，会返回字符串
- 如果请求失败，返回 `None`

示例：

```python
result = client.execute_command("get-service-info")
print(result)
```

### `ping() -> bool`

用于检查当前服务是否可达。

返回值：

- `True`: 当前请求成功
- `False`: 当前请求失败

注意：

- 这里是业务层命令检测
- 不是底层 TCP keepalive
- 也不是服务端心跳回包确认

### `get_server_version() -> str | None`

获取服务端版本号。

### `get_allowed_commands() -> list[str] | None`

获取服务端允许的命令列表。

### `is_command_allowed(command: str) -> bool`

检查某条命令是否被服务端允许。

### `get_network_config() -> dict | None`

获取服务端网络配置。

## 心跳控制

客户端暴露了一个心跳控制对象：

```python
client.heart
```

### 自动行为

- `MisakaDBClient` 初始化时会启动心跳线程
- `connect()` 时会确保心跳处于启动状态
- 已连接状态下，心跳线程会定时发送心跳帧

### 手动启动和停止

```python
client.heart.start()
client.heart.stop()
```

说明：

- `start()` 启动或恢复心跳发送
- `stop()` 停止发送心跳，但不会关闭连接

### 心跳统计字段

```python
client.heart.count
client.heart.success_count
client.heart.failure_count
client.heart.loss_rate
client.heart.running
client.heart.last_error
client.heart.last_sent_at
```

字段说明：

- `count`: 心跳发送尝试总次数
- `success_count`: 心跳发送成功次数
- `failure_count`: 心跳发送失败次数
- `loss_rate`: 当前统计口径下的失败率，单位百分比
- `running`: 当前是否处于启用状态
- `last_error`: 最近一次失败错误
- `last_sent_at`: 最近一次成功发送心跳的时间戳

### 获取完整统计

```python
stats = client.heart.stats()
print(stats)
```

返回示例：

```python
{
    "running": True,
    "interval": 10.0,
    "count": 12,
    "success_count": 12,
    "failure_count": 0,
    "loss_rate": 0.0,
    "last_error": None,
    "last_sent_at": 1710000000.123
}
```

### 关于“丢包率”的说明

当前协议没有单独的心跳响应包，因此这里的“丢包率”并不是网络层真实丢包率，而是：

```text
心跳发送失败次数 / 心跳发送尝试总次数
```

如果后续服务端增加了心跳应答包，这里的统计口径可以再升级为真正的收发确认率。

## 完整示例

```python
from apis.api import MisakaDBClient

with MisakaDBClient("127.0.0.1", 10032, heartbeat_interval=5.0) as client:
    info = client.get_service_info()
    print("服务信息:", info)

    print("服务版本:", client.get_server_version())
    print("网络配置:", client.get_network_config())
    print("允许命令:", client.get_allowed_commands())

    print("心跳是否运行:", client.heart.running)
    print("心跳次数:", client.heart.count)
    print("心跳失败率:", client.heart.loss_rate)

    client.heart.stop()
    print("手动停止心跳:", client.heart.running)

    client.heart.start()
    print("重新启动心跳:", client.heart.running)

    result = client.execute_command("get-service-info")
    print("命令结果:", result)
```



## 注意事项

- `execute_command()` 和心跳发送共用同一个 socket，并通过锁保证不会交叉写入
- 当前客户端是长连接模型，建议复用连接，不要高频短连接压测
- 如果你要做压测，优先参考 `usage/demo1.py` 的“单连接多请求”模式
- 如果服务端返回的内容不是合法 JSON，客户端会自动退回为字符串返回
