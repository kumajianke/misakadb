---
date: 2026-07-02
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
task.TaskBooks = tasktype.NewShipBuilder()
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

## 原子任务中心结构
原子中心有一个map，map管理了多个task，task管理了多个TaskBooks，也就是作业本。

原子任务中心维护了一个TasksMap，其主要目的是用于管理存储所有的任务。其map数据结构为：
```go
*xsync.MapOf[string, *Task]
```
键表示任务的唯一标志，值表示任务体指针。Task是任务体结构，包含多个作业本。
举一个例子，当我们添加一个任务的时候：
```go
center := atomic_work_center.NewAtomicWorkCenter()
ok, task_id := center.AddTask(task, 3)
```
AddTask方式会返回两个值：
- boolean 表示是否添加成功
- string 表示任务在中心的唯一标志 

![atomic](/static/atomic.png)
## TaskBooks
TaskBooks是任务的作业本，一个任务由多个作业组成，每个作业表示的是一个细分的操作，通过 TaskType 进行指定作业内容。
可见源码如下：
```go
const (
	TaskRemoveFile   TaskType = "remove_file"
	TaskModFile      TaskType = "mod_file"
	TaskRemoveFolder TaskType = "remove_folder"
)

type TaskBooks struct {
	TaskType TaskType
	Params   []string
}
```

### TaskBooks的执行和回滚

#### 执行
执行作业本的方式是采用`原子任务中心`，首先需要把`作业本`挂载到指定的Task上，然后通过addTask将Task添加到原子任务中心。添加了任务之后，就可以使用原子任务中心的如下方法进行执行了。

##### 执行的方式
- `DoNext`: 执行下一个作业;
- `DoSustain`: 从任务执行位置开始继续执行作业;

::: tip Misaka 提醒你
Task有一个属性叫做TaskCurrentIndex， 用于记录当前的任务的执行位置的，当启动服务之后这个位置不是0状态也不是over的时候表示这个任务可能在执行中途，服务器宕机了。
:::

##### 插件支持
任务中心支持用户自定义作业本的执行方式，用户可以通过实现`TaskBooks`的`TaskType`来定义自己的作业类型。你可以查看文档进行实现插件。当前所有TaskType的功能实现都是插件支持的。
