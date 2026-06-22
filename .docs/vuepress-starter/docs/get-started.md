

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

## 全局锁池

MisakaDB 在锁管理上引入了全局锁池设计，用统一的池化方式管理业务锁对象，减少重复创建与散落管理带来的开销。整个方案围绕“双池分代 + 原子引用计数 + 热点锁升级”展开，目标是在保证并发性能的同时，为后续锁回收和长期运行时的内存稳定性打下基础。

### 原理

全局锁池内部被划分为两个子池：

- `YoungPool`：存储当前高频访问、命中率较高的热点锁
- `OldPool`：存储暂时降级的锁，用于后续淘汰和回收判断

业务线程在获取锁时，会优先查询 `YoungPool`。如果未命中，再继续查询 `OldPool`。一旦在 `OldPool` 中找到目标锁，就会将该锁重新提升到 `YoungPool`，从而让热点锁始终尽量停留在更快命中的路径上。

为了降低并发争用，双池查询使用基于 `xsync.Map` 的并发 Map 实现，避免传统全局互斥锁在高并发场景下造成过多阻塞。锁对象内部还维护了一个基于原子操作的引用计数器 `RefCounter`，用于记录锁的活跃状态。每次成功加锁时增加计数，每次释放锁时减少计数，这样后续在进行分代回收时，就可以更安全地判断哪些锁已经不再被业务使用。

当前实现已经具备以下能力：

- 基于全局单例方式统一维护锁池
- 基于双池结构实现热点锁读取与旧锁提升
- 使用原子引用计数跟踪锁的使用状态
- 使用并发安全的 Map 降低锁池本身的竞争开销

当前版本中，后台 GC 守护协程与完整的“年轻代批量降级、旧生代按引用计数回收”机制仍在持续完善中。也就是说，现阶段已经完成了全局锁池的核心访问链路与数据结构设计，而分代淘汰和自动回收属于下一阶段重点增强能力。

### 设计价值

- 减少业务层重复创建锁对象的成本
- 让热点锁优先停留在高命中路径，提升锁获取效率
- 通过原子计数为后续无阻塞回收提供判断依据
- 为长期运行场景中的锁对象治理和内存回收预留扩展空间

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

本项目基于 `木兰宽松许可证` 许可证开源。详见 [LICENSE](./LICENSE) 文件。
