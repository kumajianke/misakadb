

# MisakaDB

![](./logos.png)

> [!NOTE]
> 便捷简单的文档数据库，始终如一。

`MisakaDB` 是一个轻量级的 JSON 文档数据库，支持 `JSON` 格式的内容字段存储。在 0.0.3 版本，仅支持单机模式。

## 特性

- 轻量级设计，部署简单
- 支持 JSON 格式数据存储
- Go 语言编写的高性能服务端
- Python 客户端 SDK
- 心跳保活机制
- 连接管理与命令分发

## 项目结构

```
misakadb/
├── bin/                    # 编译后的二进制文件
├── client/                 # Python 客户端
│   ├── apis/              # API 实现
│   ├── interface/        # 接口定义
│   ├── network/           # 网络通信
│   └── usage/             # 使用示例
├── clilog/                 # 日志模块
├── command/                # 命令处理
├── config/                 # 配置管理
├── engine/                 # 数据库引擎
├── misaka-doc/             # 文档
├── network/                # 网络层
├── profiles/               # 配置文件
├── safe/                   # 安全模块
├── tools/                  # 工具集
└── miusers/                # 用户管理
```

## 快速开始

### 启动服务端

#### 源码启动

需要安装 Go 语言环境后执行：

```bash
go run misaka.go [参数]
```

#### 程序启动

首先编译：

```bash
go build misaka
```

然后启动编译后的程序即可。

#### 配置文件启动

配置文件 `misaka.yaml` 内容示例：

```yaml
network:
  port: 8080
  address: 0.0.0.0
  max_conn: 100000
  retry_count: 3
  retry_delay: 30

service:
  version: "0.0.3"
```

**配置说明：**

| 配置项 | 说明 |
|--------|------|
| network.port | 服务端口 |
| network.address | 监听地址 |
| network.max_conn | 最大连接数 |
| network.retry_count | 重试次数 |
| network.retry_delay | 心跳/重试超时时间(秒) |
| service.version | 版本号 |

#### 命令行参数

```
-address    启动的地址簇
-port       启动的端口
-configs    配置文件路径【默认 misaka.yaml】
-debug      调试模式，上线前确保为 false，会在 6060 端口启动 pprof 服务
```

### 使用客户端

#### 源码运行

```bash
uv run main.py
```

#### 使用打包好的客户端

```bash
main.exe
```

#### 基础用法示例

```python
from apis import connect

# 创建并连接客户端
client = connect("127.0.0.1", 10032)

# 执行命令
result = client.execute_command("your_command")

# 关闭连接
client.close()
```

## API 参考

### MisakaDBClient

主要方法：

| 方法 | 说明 |
|------|------|
| `connect(retries, retry_delay)` | 连接到服务器 |
| `close()` | 关闭连接 |
| `get_service_info()` | 获取服务信息 |
| `execute_command(command)` | 执行命令 |
| `ping()` | 心跳检测 |
| `get_server_version()` | 获取服务器版本 |
| `is_command_allowed(command)` | 检查命令是否允许 |
| `get_network_config()` | 获取网络配置 |

### HeartbeatController

心跳控制器，用于维护与服务器的连接状态。

| 属性 | 说明 |
|------|------|
| `running` | 是否运行中 |
| `count` | 心跳总次数 |
| `success_count` | 成功次数 |
| `failure_count` | 失败次数 |
| `loss_rate` | 丢包率 |

## 许可证

本项目基于 MIT 许可证开源。详见 [LICENSE](./LICENSE) 文件。