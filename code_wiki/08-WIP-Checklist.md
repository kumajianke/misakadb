# 08. WIP 清单与启动开发指南

## 0. 开发接手前的必备认知

- 这是一个 TCP 长连接服务，不是 HTTP
- “MIQL” 当前是 `mq.<json>` 的通道协议，不是传统 SQL parser
- 引擎层已经有 `RemoveDB` 等能力，但 MIQL 的 DropDB 路由尚未接上

## 1. 本地启动清单（服务端 + 客户端）

### 1.1 初始化 profiles（必须）

首次启动前必须生成密钥与用户表，否则服务端会因为缺失 `profiles/master.mikey` 直接 panic：

```bash
go run ./tools/misaka-tools.go sys-init
```

然后立即修改 root 密码（需要在本机终端中输入密码；CI/非 TTY 环境无法安全读取）：

```bash
go run ./tools/misaka-tools.go chpwd root
```

### 1.2 启动服务端

源码启动（默认会加载 `./profiles/misaka.yaml`）：

```bash
go run ./misaka.go
```

常用参数（覆盖配置）：

```bash
go run ./misaka.go -address 0.0.0.0 -port 10032 -debug=false
```

### 1.3 启动 Python 客户端 CLI

```bash
cd client
python main.py --address 127.0.0.1 --port 10032 --mode shell
```

CLI 会先调用 `get-service-info`，随后要求登录后才允许执行 `mq.` 命令。

## 2. 当前 WIP 状态清单（按风险优先级）

### 2.1 DropDB 逻辑未实现（高优先级）

现状：

- MIQL 路由已经支持 `drp-dat` 分支：[RunMson](file:///workspace/engine/share/miql.go#L14-L19)
- 客户端也能生成 `{"active":"drp-dat","name":"xxx"}`：[MiQL.dropDB](file:///workspace/client/mql/MQ.py#L23-L26)
- 但服务端 `MiqlDropDB` 仅做 role 校验后就 `return nil`：[MiqlDropDB](file:///workspace/engine/share/miqlLogic.go#L40-L51)

### 2.2 引擎缓存回填 bug（中高）

`GetDBEngine` 回填缓存时使用了固定 key `"dbname"`，导致缓存无法按 DB 命中：

- [engine_factory.go](file:///workspace/engine/dispatch/engine_factory.go#L28-L31)

### 2.3 TinyDB 的 loader 未完成（中）

- `WriteLoader` / `ReadLoader` 当前为空实现：[TinyDBLoaderImp](file:///workspace/engine/tinydb/components/TinyDBLoader.go#L44-L57)

### 2.4 DML 执行器未接入（中）

接口已定义 `MiQLExecutorCore`（Insert/Delete/Update/Search），但当前主链路未出现调用位置：

- [MiQLExecutorCore](file:///workspace/engine/base/BaseCore.go#L74-L79)

### 2.5 协议一致性与错误语义（中）

- 服务端响应前缀存在 `[ok]`、`[err]`、`[error]` 三套风格，客户端对三者处理略有差异
- 建议明确：
  - `[ok]`：成功
  - `[err]`：业务失败（可预期）
  - `[error]`：协议/输入错误（不可预期或安全拦截）

## 3. 下一步开发任务拆解：实现 DropDB（DROPDB）

### 3.1 目标行为（建议作为验收标准）

- 输入：登录后的 MIQL 命令 `mq.{"active":"drp-dat","name":"student"}`
- 权限：仅允许 `root` 角色执行
- 成功：删除 `db-datas/student/`，并返回 `[ok]drop db is ok!`
- 失败：
  - 未登录：由 dispatch 层拦截并返回 `[error]you must login first`
  - 非 root：返回 `[err]<没有权限>`
  - DB 不存在：返回 `[err]cannot found this db: <name>`（或更明确的错误）

### 3.2 实现步骤（建议落地顺序）

1. **完善权限失败的中断逻辑**
   - `VerifyRole` 失败后必须 `return err`，避免继续执行并最终 `return nil`
2. **DB 名校验（安全）**
   - 限制 dbname 只能为 `[a-zA-Z0-9_-]`（至少禁止 `..`、路径分隔符），避免路径穿越影响 `os.RemoveAll`
3. **选择引擎并执行删除**
   - 如果 drop 请求没有指定 engine：
     - 使用 `engine_dispatch.GetDBEngine(dbname)` 从 meta 推断引擎
   - 调用 `dbEngine.RemoveDB(dbname)`
4. **清理引擎缓存映射**
   - 从 `RegisterCenter.MapperDBEngine` 中删除该 dbname 的映射（并修复缓存回填 bug）
5. **返回统一响应**
   - 成功：`[ok]...`
   - 失败：`[err]...`（保证客户端 CLI 能清晰展示）

### 3.3 建议补充的测试/验证

- 手工冒烟：
  - `mq.createDB("student")` → 创建成功
  - `mq.dropDB("student")` → 删除成功
  - 再次 `mq.dropDB("student")` → 返回“DB 不存在”
- 并发安全：
  - 在创建/删除同名 DB 时不应出现半创建状态（需依赖行锁粒度是否足够）

## 4. 推荐的工程化增强（可选，但利于后续扩展）

- 将 MIQL 从 “Mson JSON” 演进为可扩展的 AST/命令结构（保留向后兼容）
- 为 `ServiceConnContext` 增加更明确的会话状态（例如 current db、tx context）
- 为协议增加心跳响应包（客户端统计丢包率可变为真实确认率）
- 将 `db-datas` 根目录改为可配置（已有 `private.storage.path` 字段，但尚未被引擎真正使用）

