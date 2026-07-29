package pluginsx

import (
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"sort"
	"strings"
	"sync"

	"github.com/puzpuzpuz/xsync/v3"
)

var (
	pluginsX     PluginsBus
	pluginsXOnce sync.Once
)

type PluginsBus struct {
	TaskTypes       []tasktype.TaskType                    // misaka所加载的任务类型
	TaskTypeAlias   xsync.MapOf[string, tasktype.TaskType] // 任务别名到 TaskType 的映射
	TaskTypeActions xsync.MapOf[tasktype.TaskType, any]    // 不同对应的任务类型的操作句柄
	TaskTypeRoll    xsync.MapOf[tasktype.TaskType, any]    // 不同对应的任务类型的回滚句柄
	TaskAliasDocs   xsync.MapOf[string, TaskAliasDoc]      // 别名说明文档
}

type TaskAliasDoc struct {
	Plugin   string
	Alias    string
	Desc     string
	TaskType tasktype.TaskType
}

func GetPluginsBus() *PluginsBus {
	pluginsXOnce.Do(func() {
		pluginsX = PluginsBus{
			TaskTypes:       make([]tasktype.TaskType, 0),
			TaskTypeAlias:   *xsync.NewMapOf[string, tasktype.TaskType](),
			TaskTypeActions: *xsync.NewMapOf[tasktype.TaskType, any](),
			TaskTypeRoll:    *xsync.NewMapOf[tasktype.TaskType, any](),
			TaskAliasDocs:   *xsync.NewMapOf[string, TaskAliasDoc](),
		}
	})
	return &pluginsX
}

func (bus *PluginsBus) StoreTaskAliasDoc(doc TaskAliasDoc) {
	if strings.TrimSpace(doc.Alias) == "" {
		return
	}
	bus.TaskAliasDocs.Store(doc.Alias, doc)
}

func (bus *PluginsBus) StoreTaskTypeAlias(alias string, taskType tasktype.TaskType) {
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return
	}
	bus.TaskTypeAlias.Store(alias, taskType)
}

func (bus *PluginsBus) ResolveTaskType(alias string) (tasktype.TaskType, bool) {
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return "", false
	}
	return bus.TaskTypeAlias.Load(alias)
}

func (bus *PluginsBus) GetTaskAliasDocsSnapshot() []TaskAliasDoc {
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
