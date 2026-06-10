# 03. 模块职责

## clilog（日志）

- 目的：统一带颜色的日志输出（INFO/SUCCESS/WARNING/ERROR）
- 位置：[logger.go](file:///workspace/clilog/logger.go)

## config（配置）

- `MisakaConfigure`：yaml/json 映射结构
- `InitGlobalMisakaConfigure`：通过 `sync.Once` 初始化全局配置
- 注意：网络层收包依赖 `GetGlobalNetworkConfigure()`，因此必须在启动网络服务前先完成 init

相关代码：

- 配置结构：[MisakaConfigure.go](file:///workspace/config/MisakaConfigure.go)
- 全局加载与默认值：[mapper.go](file:///workspace/config/mapper.go)

## network（网络与连接管理）

### network/core

- `ServiceCore.Run()`：监听 TCP，Accept 新连接
- `handlerConn()`：每连接 goroutine 循环 `Recv → Dispatch`

位置：[ServiceCore.go](file:///workspace/network/core/ServiceCore.go)

### network/context

- `ServiceConnContext`：连接级状态
- `Recv()`：封装 `RecvWithHeart`
- `Send()`：`len + payload` 的服务端响应写包（无 type byte）

位置：[ServiceConnContext.go](file:///workspace/network/context/ServiceConnContext.go)

### network/SockShare

- `RecvWithHeart()`：读 `4字节长度 + 1字节类型 + payload`
  - `type=0x01` 视为心跳包直接忽略并继续等待下一个数据包
  - 内置最大载荷限制 `16MB` 防止 OOM

位置：[Heater.go](file:///workspace/network/SockShare/Heater.go)

### network/RegisterCenter

- 全局单例：维护 `ConnectQueue`、`MasterKey`、`MapperDBEngine`
- 连接队列：通过 channel 控制最大连接数；连接释放在 `handlerConn` defer 中处理

位置：[core.go](file:///workspace/network/RegisterCenter/core.go)

### network/Onces

- `SafeConn`：用 `sync.Once` 确保连接只关闭一次，避免并发 close panic

位置：[ServiceOnce.go](file:///workspace/network/Onces/ServiceOnce.go)

## command（命令系统）

### CommandDispatch

- 统一入口：`MiqlCommDispatch.Dispatch(ctx, command)`
- 分两类：
  - `mq.` 前缀：进入 MIQL 通道（要求已登录）
  - 否则：按结构体字段 tag `mapper:"xxx"` 反射路由到具体 handler

位置：[CommandDispatch.go](file:///workspace/command/CommandDispatch.go)

### Func（内置普通命令）

- `get-service-info`：输出配置（可通过 `service.hide_info` 关闭该能力）
- `login <user> <pass>`：登录并写入 `ctx.LoginUser`
- `exit`：关闭连接

位置：[Func.go](file:///workspace/command/Func.go)

## engine（引擎层）

### engine/Mson

- `MsonParse`：MIQL payload 的 JSON 结构体
- `NewMsonParse()`：服务端将 `mq.` 后内容直接 `json.Unmarshal`

位置：[MsonParse.go](file:///workspace/engine/Mson/MsonParse.go)

### engine/share（MIQL 路由与逻辑）

- `RunMson()`：按 `active` 分发（当前支持 `cre-dat`、`drp-dat`）
- `MiqlCreateDB()`：`NewEngine` → `InitLoader` → 返回 `[ok]`
- `MiqlDropDB()`：目前仅做 role 校验，属于 WIP（详见 08-WIP-Checklist）

位置：

- 路由：[miql.go](file:///workspace/engine/share/miql.go)
- 逻辑：[miqlLogic.go](file:///workspace/engine/share/miqlLogic.go)

### engine/base（接口与公共能力）

- `BaseEngineCore`：统一引擎接口（锁/loader/baker/executor/remove）
- `EngineLockerSupport`：基于全局锁池实现“引擎级锁 + 行级锁”
- `ShareLoaderDBMetaName()`：从 `db-datas/<db>/meta.json` 读取引擎名

位置：[BaseCore.go](file:///workspace/engine/base/BaseCore.go)

### engine/tinydb（当前唯一引擎）

- `TinyDBCore.RemoveDB()`：执行 `os.RemoveAll("./db-datas/<dbname>")` 删除 DB
- `TinyDBLoaderImp.InitLoader()`：创建 DB 目录结构并写 `meta.json`
- 注意：`WriteLoader/ReadLoader` 目前为空实现

位置：

- 引擎核心：[TinyCore.go](file:///workspace/engine/tinydb/TinyCore.go)
- loader：[TinyDBLoader.go](file:///workspace/engine/tinydb/components/TinyDBLoader.go)

## miusers + safe（用户与安全）

- `safe`：AES-GCM 加解密；`InitPassword()` 初始化 `profiles/master.mikey`
- `miusers`：用户文件 `profiles/user.dat`（加密存储），使用 bcrypt 保存密码 hash

位置：

- AES：[AES256.go](file:///workspace/safe/AES256.go)
- 用户管理：[userManager.go](file:///workspace/miusers/userManager.go)

## tools（运维/开发工具）

- `misaka-tools`：命令行工具，支持初始化密钥/用户、添加用户、修改密码、改角色等

入口：[misaka-tools.go](file:///workspace/tools/misaka-tools.go)  
命令实现：[CommandExecute.go](file:///workspace/command/ToolsCommands/CommandExecute.go)

