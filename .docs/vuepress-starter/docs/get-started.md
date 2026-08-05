

# MisakaDB

![](./logos.png)

> 便捷简单的文档数据库，始终如一。

::: tip
misaka的命名来自动漫《魔法禁书目录》、《某科学的超电池炮》的炮姐，御坂美琴中的御坂(misaka)及衍生的御坂网络，在后期我们也会仿造misaka来实现对应的分布式技术。
:::
`MisakaDB` 是一个轻量级的 JSON 文档数据库，支持 `JSON` 格式的内容字段存储。在 0.0.3 版本，仅支持单机模式。在0.17版本之后，支持插件模式，开发者可以在Misaka源码中开发自定义的插件，从而提高CRUD的效率。

## 特性

- 轻量级设计，部署简单
- 支持 JSON 格式数据存储
- Go 语言编写的高性能服务端
- Python 客户端 SDK
- 心跳保活机制
- 连接管理与命令分发

## 全局锁池

MisakaDB 在锁管理上引入了全局锁池设计，用统一的池化方式管理业务锁对象，减少重复创建与散落管理带来的开销。整个方案围绕“双池分代 + 原子引用计数 + 热点锁升级”展开，目标是在保证并发性能的同时，为后续锁回收和长期运行时的内存稳定性打下基础。


### 设计价值

- 减少业务层重复创建锁对象的成本
- 让热点锁优先停留在高命中路径，提升锁获取效率
- 通过原子计数为后续无阻塞回收提供判断依据
- 为长期运行场景中的锁对象治理和内存回收预留扩展空间

## 项目结构

```
misaka_db/
├── .docs/                  # VuePress 文档工程
│   └── vuepress-starter/   # 文档站点源码与构建配置
├── atomic/                 # 原子任务、事件总线与文件操作能力
│   ├── EventBus/           # 原子任务事件总线
│   ├── atomicFileHandler/  # 原子文件写入等操作
│   └── atomicWorkCenter/   # 原子任务中心、序列化与任务类型
├── bin/                    # 编译产物与运行时配置
├── client/                 # Python 客户端 SDK 与打包脚本
│   ├── apis/               # 客户端 API
│   ├── interface/          # 状态与接口定义
│   ├── mql/                # MQL 相关调用封装
│   ├── network/            # 客户端网络通信
│   └── usage/              # 使用示例
├── clilog/                 # 日志输出模块
├── command/                # 命令分发与工具命令实现
├── config/                 # 配置加载与映射
├── db-datas/               # 本地数据库数据目录
├── engine/                 # 存储引擎与查询实现
│   ├── Mson/               # MSON 解析
│   ├── base/               # 引擎基础能力
│   ├── db_list/            # 数据库列表管理
│   ├── dispatch/           # 引擎分发工厂
│   ├── share/              # MIQL 共享逻辑
│   └── tinydb/             # TinyDB 核心、索引与数据集
├── lock/                   # 全局锁池与锁前缀
├── miusers/                # 用户管理
├── network/                # 服务端网络层与注册中心
├── profiles/               # 默认配置与用户数据
├── safe/                   # 加密与安全模块
├── shares/                 # 跨模块共享工具
├── tools/                  # 工具程序入口
├── tui/                    # 终端 UI 能力
├── misaka.go               # 服务主入口
└── go.mod                  # Go 模块定义
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
