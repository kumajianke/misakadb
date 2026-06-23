---
date: 2026-06-22
title: 全局锁池
category:
  - 内部实现
tag:
  - 全局锁池
description: 全局锁池的实现
---

# 全局锁池

::: tip 什么是全局锁池
全局锁池可以帮助我们快速的在项目中对某个资源、操作句柄进行上锁。
:::

## 全局锁池的实现
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



## 使用方法
```go
import (
	"misakadb/lock/global_lock"
)

lock , unlock, err := global_lock.GetOrStoreGlobalLock("lock_key")
defer unlock()
// 独占锁之后需要进行的操作
```

::: warning
我们获取的第一个返回值虽然是lock，但是不要使用
```go
defer lock.Unlock()
```
来释放锁，这样会导致锁无法被全局锁池释放!
:::

## 知其所以然
有兴趣的同学可以去看看我们的博客: [内部实现原理](http://blog.6ugo.cn/?view=tech&id=316efd3b-c398-4fb2-bca0-3655893a0934)

## 源码解析
源码在代码中的 `lock/global_lock/GlobalLock.go` 文件中。

```go
// 本代码引用的是0.17的源码 如果有出入还请对照参考
...
var once sync.Once
var globalLocksPool *GlobalLocksPool

// 全局锁的结构
type GlobalLocks struct {
  ...
}

// 全局锁池的结构
type GlobalLocksPool struct {
	...
}

// 获取一个全局单例，用于存储所有的锁
func GetGlobalLockPool() *GlobalLocksPool {
  ...
}

// 获取年轻代的锁的引用快照
func (p *GlobalLocksPool) GetYoungPoolSnapshot() map[string]int32 {
	...
}

// 获取旧代的锁的引用快照
func (p *GlobalLocksPool) GetOldPoolSnapshot() map[string]int32 {
	...
}

// 获取一个锁，如果锁是第一次创建，自动注册到全局锁池
func GetOrStoreGlobalLock(lock_name string, lock_method string) (*GlobalLocks, func(), error) {
	...
}

func lockPoolsGCThread() {
	//回收机制触发
}

// 启动全局锁池的回收线程
func StartLockPoolGC() {
	...
}

```