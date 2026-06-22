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

## 什么是全局锁池
::: tip
全局锁池可以帮助我们快速的在项目中对某个资源、操作句柄进行上锁。
:::

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