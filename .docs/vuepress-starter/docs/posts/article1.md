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

lock , unlock, err :=global_lock.GetOrStoreGlobalLock("lock_key")
```
