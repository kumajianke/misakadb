# 01. 项目概览

## 项目定位

MisakaDB 是一个轻量级 JSON 文档数据库：

- 服务端：Go 编写，提供 TCP 长连接服务与命令分发
- 客户端：提供 Python SDK + CLI
- 存储：当前默认引擎为 `tinydb`，使用文件系统目录组织数据库数据

## 核心概念

### Command vs. MIQL

- 普通命令：以空格分隔的纯文本命令（例如 `login <u> <p>`、`get-service-info`）
- MIQL：以 `mq.` 前缀标识的“数据库语句”通道；当前实现并不是传统 SQL，而是：
  - 客户端将一个 JSON（称为 `Mson`）序列化后拼到 `mq.` 后面
  - 服务端收到后直接 JSON 反序列化为 `MsonParse` 结构体，再按 `active` 字段路由

### Mson（MIQL Payload）

服务端结构体定义为：

- `active`：动作标识（如 `cre-dat` / `drp-dat`）
- `name`：数据库名
- `engine`：引擎名（如 `tinydb`）

对应代码：[MsonParse](file:///workspace/engine/Mson/MsonParse.go)

### 引擎（Engine）

每个数据库由“引擎核心”对象承载，统一通过接口抽象（锁、loader、备份、执行器、删除 DB 等能力）：

- 接口定义：[BaseEngineCore](file:///workspace/engine/base/BaseCore.go#L84-L90)
- 工厂创建：[engine_factory.go](file:///workspace/engine/dispatch/engine_factory.go)

### RegisterCenter

一个全局单例，承担：

- 连接队列（最大连接数控制）
- master key（AES 密钥）读取与共享
- DB 名 → 引擎名的缓存映射（`MapperDBEngine`）

实现位置：[RegisterCenter](file:///workspace/network/RegisterCenter/core.go)

### 全局锁池（GlobalLocksPool）

用于统一管理业务锁对象，避免到处散落创建锁，当前实现为“young/old 双池 + 引用计数 + 定期 GC”：

- 关键实现：[GlobalLock.go](file:///workspace/lock/global_lock/GlobalLock.go)

## 仓库结构（按用途分组）

### Go 服务端（核心）

- `misaka.go`：服务端入口
- `network/`：TCP 服务、连接上下文、心跳收包、注册中心
- `command/`：普通命令与 MIQL 分发
- `engine/`：引擎层（base 接口、tinydb 实现、mson 解析、miql 路由）
- `miusers/`：用户管理（bcrypt + AES 加密文件）
- `safe/`：AES-GCM 加解密与密钥初始化
- `lock/`：全局锁池
- `clilog/`：日志输出

### 工具（Go）

- `tools/misaka-tools.go` + `command/ToolsCommands/`：初始化密钥/用户、管理用户等（面向运维/开发）

### Python 客户端

- `client/apis/api.py`：Python SDK（MisakaDBClient + HeartbeatController）
- `client/main.py`：CLI 入口
- `client/network/sock.py`：底层 socket 读写协议实现
- `client/mql/MQ.py`：MIQL 语句构造器（生成 `mq.<json>`）

## 运行所需的关键文件

服务端启动依赖 `profiles/` 下的敏感文件：

- `profiles/master.mikey`：AES 32 字节密钥（缺失会导致服务端直接 panic）
- `profiles/user.dat`：加密后的用户表（login / role 校验依赖它）
- `profiles/misaka.yaml`：默认配置文件（misaka.go 缺省会读取这个路径）

