---
date: 2026-06-22
title: 原子任务中心
category:
  - 内部实现
tag:
  - 原子任务
description: 原子任务中心介绍
---

# 原子任务中心的作用
当客户端发来<span style="color:#0099FF">`miql`</span>之后，misakadb要统一进行解析，为了保持miql的任务的原子性、一致性、隔离性、持久性这四个属性，misakaDB需要一个原子任务中心来管理所有的任务。

> 如当用户发送一个删除任务，这个任务拆分下来就是需要对文件、JSON等进行删除的多个子任务。这些任务需要保证要么不执行要么全部执行。

所以我们执行多个需要多步骤的任务不能直接进行执行，而是把子任务一次性追加到我们的新任务中。
```go
atomic := atomic_work_center.NewAtomicWorkCenter();
// 创建一个原子任务中心
task = atomic_work_center.NewTask(nil);
// 创建一个空任务
task.TaskBody = tasktype.NewShipBuilder()
	.Add(tasktype.TaskModFile, "xxx.filename.txt")
	.Build();
// 构建一个任务体，任务体包含需要原子执行的多个任务子作业

ok, task_id := atomic.addTask(task, 3) // 尝试添加任务到原子中心 可以重试3次
// ok: 表示是否添加成功
// task_id: 表示任务在中心的唯一标志
```

# 原子日志序列器
<span style="color:#0099FF">**AtomicWorkCenter**</span> 管理着所有的任务，我们将这些任务数据可以称为WAL日志。防止宕机等问题导致的数据丢失，我们会在启动之后启动对应的序列器协程。如果需要让协程启动序列化任务可以使用AtomicWorkEventBus总线进行信号通知。

```go
var eb *eventbus.AtomicWorkEventBus
eb = eventbus.NewAtomicWorkCenterEventBus()
eb.EventBus <- "sync-to-local" // 通知序列器协程启动序列化任务
```

我们也可以加载本地WAL日志到内存:

```go
workCenterSerializer := atomic_work_center.LoadWorkCenterSerializer()
```

WAL的日志默认储存路径是：`.data/work_center.json` ， 如果日志大于1MB的时候，系统会对其进行分片存储。atomic_work_center的本质就是 WAL日志的记录。