package pluginsx

import (
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	pluginsXInterface "misakadb/plugins/pluginsx/pluginsx_interface"
	"sort"
	"strings"
	"sync"

	"github.com/puzpuzpuz/xsync/v3"
)

var (
	pluginsX     PluginsX
	pluginsXOnce sync.Once
)

type PluginsX struct {
	TaskTypes       []tasktype.TaskType                                                  // misaka所加载的任务类型
	TaskTypeAlias   xsync.MapOf[string, tasktype.TaskType]                               // 任务别名到 TaskType 的映射
	TaskTypeActions xsync.MapOf[tasktype.TaskType, pluginsXInterface.FuncTaskTypeAction] // 不同对应的任务类型的操作句柄
	TaskTypeRoll    xsync.MapOf[tasktype.TaskType, pluginsXInterface.FuncTaskTypeRoll]   // 不同对应的任务类型的回滚句柄
	TaskAliasDocs   xsync.MapOf[string, TaskAliasDoc]                                    // 别名说明文档
	TaskCombo       xsync.MapOf[string, pluginsXInterface.FuncTaskCombo]                 // 不同对应的混招
	AfterTasks      xsync.MapOf[tasktype.TaskType, pluginsXInterface.AfterTask]          // 任务结束之后需要执行的内容
}

type TaskAliasDoc struct {
	Plugin   string
	Alias    string
	Desc     string
	TaskType tasktype.TaskType
}

/*
*
单例获取plugins的
*/
func GetPluginsX() *PluginsX {
	pluginsXOnce.Do(func() {
		pluginsX = PluginsX{
			TaskTypes:       make([]tasktype.TaskType, 0),
			TaskTypeAlias:   *xsync.NewMapOf[string, tasktype.TaskType](),
			TaskTypeActions: *xsync.NewMapOf[tasktype.TaskType, pluginsXInterface.FuncTaskTypeAction](),
			TaskTypeRoll:    *xsync.NewMapOf[tasktype.TaskType, pluginsXInterface.FuncTaskTypeRoll](),
			TaskAliasDocs:   *xsync.NewMapOf[string, TaskAliasDoc](),
			TaskCombo:       *xsync.NewMapOf[string, pluginsXInterface.FuncTaskCombo](),
			AfterTasks:      *xsync.NewMapOf[tasktype.TaskType, pluginsXInterface.AfterTask](),
		}
	})
	return &pluginsX
}

func (bus *PluginsX) StoreTaskAliasDoc(doc TaskAliasDoc) {
	if strings.TrimSpace(doc.Alias) == "" {
		return
	}
	bus.TaskAliasDocs.Store(doc.Alias, doc)
}

func (bus *PluginsX) StoreTaskTypeAlias(alias string, taskType tasktype.TaskType) {
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return
	}
	bus.TaskTypeAlias.Store(alias, taskType)
}

func (bus *PluginsX) ResolveTaskType(alias string) (tasktype.TaskType, bool) {
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return "", false
	}
	return bus.TaskTypeAlias.Load(alias)
}

func (bus *PluginsX) GetTaskAliasDocsSnapshot() []TaskAliasDoc {
	docs := make([]TaskAliasDoc, 0)
	bus.TaskAliasDocs.Range(func(_ string, doc TaskAliasDoc) bool {
		docs = append(docs, doc)
		return true
	})

	sort.Slice(docs, func(i, j int) bool {
		if docs[i].Plugin == docs[j].Plugin {
			return docs[i].Alias < docs[j].Alias
		}
		return docs[i].Plugin < docs[j].Plugin
	})
	return docs
}
