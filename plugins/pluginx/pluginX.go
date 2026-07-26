package pluginsx

import (
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"sync"

	"github.com/puzpuzpuz/xsync/v3"
)

var (
	pluginsX     PluginsBus
	pluginsXOnce sync.Once
)

type PluginsBus struct {
	TaskTypes       []tasktype.TaskType                 // misaka所加载的任务类型
	TaskTypeActions xsync.MapOf[tasktype.TaskType, any] // 不同对应的任务类型的操作句柄
	TaskTypeRoll    xsync.MapOf[tasktype.TaskType, any] //不同对应的任务类型的回滚句柄
}

func GetPluginsBus() *PluginsBus {
	pluginsXOnce.Do(func() {
		pluginsX = PluginsBus{
			TaskTypes:       make([]tasktype.TaskType, 0),
			TaskTypeActions: *xsync.NewMapOf[tasktype.TaskType, any](),
			TaskTypeRoll:    *xsync.NewMapOf[tasktype.TaskType, any](),
		}
	})
	return &pluginsX
}
