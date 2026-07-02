package pluginsloader

import tasktype "misakadb/atomic/atomicWorkCenter/TaskType"

var (
	PluginsTaskTypes              = make([]tasktype.TaskType, 0)       // 通过插件加载的所有的任务类型列表
	PluginsTaskTypeMethodRunnable = make(map[tasktype.TaskType]func()) // 通过插件加载的任务类型所有的可执行方法
	PluginsTaskTypeMethodRollBack = make(map[tasktype.TaskType]func()) // 通过插件加载的任务类型所有的回滚方法
)
