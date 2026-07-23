package pluginsx

import (
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"sync"
)

var pluginsX PluginsBus

type PluginsBus struct {
	AllTaskType []tasktype.TaskType // misaka所加载的任务类型
}

func GetPluginsBus() *PluginsBus {
	var once sync.Once
	once.Do(func() {
		pluginsX = PluginsBus{}
	})
	return &pluginsX
}
