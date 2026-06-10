# 06. 引擎层与 TinyDB

## 引擎接口（BaseEngineCore）

核心接口定义在：[BaseEngineCore](file:///workspace/engine/base/BaseCore.go#L84-L90)

- 锁能力：
  - `Lock()`：引擎级锁
  - `GetRowLock(name)`：行级锁（当前通常用来锁 DB 名）
- 组件能力：
  - `DBLoader()`：初始化/读写落盘组件
  - `DBBaker()`：备份组件（目前未见完整实现）
  - `MiQLExecutor()`：DML 执行器（Insert/Delete/Update/Search，当前未接入主链路）
- 管理能力：
  - `RemoveDB(dbname)`：删除数据库（DropDB 应调用这里）

## 引擎工厂与缓存

### NewEngine

当前仅支持：

- `tinydb`

实现：[NewEngine](file:///workspace/engine/dispatch/engine_factory.go#L9-L15)

### GetDBEngine（从 meta 推断引擎）

`GetDBEngine(dbname)` 的设计意图是：

1. 尝试从 `RegisterCenter.MapperDBEngine` 缓存命中
2. 读 `db-datas/<db>/meta.json` 得到引擎名（`ShareLoaderDBMetaName`）
3. 回填缓存并返回引擎实例

实现：[GetDBEngine](file:///workspace/engine/dispatch/engine_factory.go#L17-L32)

注意：当前回填缓存处使用了固定 key `"dbname"` 而不是 `dbname`，这会导致缓存逻辑失效（WIP 修复点之一）。

## TinyDB 的目录结构

创建 DB 时，会创建：

```
db-datas/<dbName>/
  .db/
    meta.json
```

创建代码：[TinyDBLoaderImp.InitLoader](file:///workspace/engine/tinydb/components/TinyDBLoader.go#L58-L124)

### meta.json 内容

由 `filejson.NewTinyDBMeta` 生成（结构体定义在 `engine/tinydb/FileJson` 下），写入位置：

- `./db-datas/<dbName>/.db/meta.json`

## CreateDB（InitLoader）

`InitLoader` 关键步骤：

1. 校验 DB 是否已存在（`os.Stat(db-datas/<name>)`）
2. 创建 DB 根目录（`0700`）
3. 创建 `.db` 目录（`0700`）
4. 写入 `meta.json`（`0600`）
5. Windows 平台下尝试隐藏 `.db` 目录（`attrib +h`）

实现：[TinyDBLoaderImp.InitLoader](file:///workspace/engine/tinydb/components/TinyDBLoader.go#L58-L124)

## DropDB（RemoveDB）

TinyDB 已实现删除能力：

- 使用“行级锁”锁住 `this.Name`
- `os.RemoveAll("./db-datas/<dbname>")`

实现：[TinyDBCore.RemoveDB](file:///workspace/engine/tinydb/TinyCore.go#L23-L43)

目前缺口不在引擎实现，而在 MIQL 的 `drp-dat` 逻辑未调用 `RemoveDB`（详见 08-WIP-Checklist）。

