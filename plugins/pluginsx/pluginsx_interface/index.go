package pluginsxInterface

import tasktype "misakadb/atomic/atomicWorkCenter/TaskType"

type FuncTaskTypeAdd = func(allTaskTypelst []tasktype.TaskType) ([]tasktype.TaskType, error)
type FuncTaskTypeAction = func(taskType tasktype.TaskType, params []string) error
type FuncTaskTypeRoll = func(taskType tasktype.TaskType, params []string) error
type FuncTaskCombo = func(params []string) (error, string)
