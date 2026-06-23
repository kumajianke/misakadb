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
当客户端发来miql之后，misakadb要统一进行解析，为了保持miql的任务的原子性、一致性、隔离性、持久性这四个属性，misakaDB需要一个原子任务中心来管理所有的任务。

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