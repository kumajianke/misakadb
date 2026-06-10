# 04. 网络协议与心跳

## 总览

MisakaDB 当前使用 TCP 长连接协议：

- 客户端（Python）主动发送心跳包
- 服务端收包时识别心跳并忽略
- 服务端响应采用更简单的 framing（不带 type byte）

## Client → Server 帧格式

由 Python 客户端实现：[client/network/sock.py](file:///workspace/client/network/sock.py#L50-L67)

```
4 bytes  length (big-endian, N)
1 byte   msg_type
N bytes  payload
```

### msg_type

- `0x00`：普通请求（payload 为 UTF-8 命令字符串）
- `0x01`：心跳包（payload 为空字节串）

## Server 收包（带心跳）

服务端实现：[RecvWithHeart](file:///workspace/network/SockShare/Heater.go)

- 设置连接 deadline：`now + retryDelay`
- 连续尝试读取 4 字节长度
  - 如果超时，累加 `errorRecvCounter`
  - 超过 `retryCount` 后返回错误（服务端会关闭连接）
- 读取 1 字节 type
  - `0x01`：心跳包，记录日志并继续等待下一帧
  - 其他：读取 `N` 字节 payload 作为命令字符串

载荷限制：

- 内置 `MaxPayloadSize = 16MB`，避免客户端传入超大 length 导致服务端 OOM

## Server → Client 帧格式

服务端发送实现：[ServiceConnContext.Send](file:///workspace/network/context/ServiceConnContext.go#L43-L56)

```
4 bytes  length (big-endian, N)
N bytes  payload
```

注意：服务端响应不带 `msg_type` 字段。

客户端接收实现：[client/network/sock.py recv_bytes](file:///workspace/client/network/sock.py#L68-L78)

## 超时与重试语义

服务端侧的“心跳重试次数/延迟”来自配置：

- `network.retry_count`
- `network.retry_delay`

读取位置：[GetGlobalNetworkConfigure](file:///workspace/config/mapper.go#L65-L70)

这两个值会直接影响 `RecvWithHeart` 的 deadline 与超时重试次数，从而决定连接被判定为“超时断开”的窗口。

