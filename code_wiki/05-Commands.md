# 05. 命令体系与分发

## 命令分两层

1. 普通命令：纯文本（空格分隔），由服务端 `command` 模块通过反射 mapper 分发
2. MIQL 通道：以 `mq.` 开头，后面直接跟 JSON 字符串（Mson）

分发入口：[MiqlCommDispatch.Dispatch](file:///workspace/command/CommandDispatch.go#L28-L85)

## 普通命令（Command）

### get-service-info

- 作用：返回当前服务配置的 JSON（前缀 `[ok]`），供客户端面板展示
- 代码：[ImpGetServiceInfo](file:///workspace/command/Func.go#L13-L52)
- 受控开关：`service.hide_info`

### login <username> <password>

- 作用：鉴权并把 `ctx.LoginUser = username`
- 代码：[ImpLogin](file:///workspace/command/Func.go#L64-L87)
- 密码校验：bcrypt hash 对比（用户表来自 `profiles/user.dat`）

### exit

- 作用：关闭连接
- 代码：[ImpExit](file:///workspace/command/Func.go#L54-L62)

## MIQL（mq.<json>）

### MIQL 的实际形态

当前 MIQL 并不是 SQL 字符串，而是“客户端拼接的 JSON 语句”：

- Python 构造器：[client/mql/MQ.py](file:///workspace/client/mql/MQ.py)
- 服务端解析：[NewMsonParse](file:///workspace/engine/Mson/MsonParse.go#L14-L24)

### 登录前置条件

服务端要求 MIQL 必须在登录后执行：

- 未登录直接返回：`[error]you must login first`
- 检查逻辑：[CommandDispatch.go](file:///workspace/command/CommandDispatch.go#L33-L38)

## MIQL 动作路由（active）

服务端按 `msonParse.Active` 路由：

- `cre-dat`：创建 DB
- `drp-dat`：删除 DB（WIP）

路由位置：[RunMson](file:///workspace/engine/share/miql.go#L9-L22)

## 已实现：CreateDB

### 客户端语句示例

Python：

```python
from mql.MQ import MiQL
mq = MiQL(cli)
res = mq.createDB("student", engine="tinydb").shot()
```

等价的 MIQL wire 命令（概念上）：

```text
mq.{"active":"cre-dat","name":"student","engine":"tinydb"}
```

### 服务端执行路径

1. `NewEngine(engineName, dbName)`
2. `DBLoader().InitLoader(mson)`
3. `Send("[ok]create db is ok!")`

实现：[MiqlCreateDB](file:///workspace/engine/share/miqlLogic.go#L12-L38)

## WIP：DropDB

客户端已有 `dropDB(name)` 构造（`active=drp-dat`），但服务端逻辑未落地：

- 客户端构造：[MiQL.dropDB](file:///workspace/client/mql/MQ.py#L23-L26)
- 服务端逻辑占位：[MiqlDropDB](file:///workspace/engine/share/miqlLogic.go#L40-L51)

下一步应当把“权限校验 → 选择引擎 → 调用 RemoveDB → 清理缓存 → 返回 [ok]/[err]”串起来，详见 08-WIP-Checklist。

