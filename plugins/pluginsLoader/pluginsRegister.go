package pluginsloader

import (
        "fmt"
        tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
        "misakadb/config"
        pluginsx "misakadb/plugins/pluginx"
        pluginxInterface "misakadb/plugins/pluginx/pluginx_interface"
        "sort"
        "strings"
        "sync"
)

var (
	loadedPlugins sync.Map
	taskTypeAlias sync.Map
        taskAliasDocs sync.Map
)

type TaskAliasDoc struct {
        Plugin   string
        Alias    string
        Desc     string
        TaskType tasktype.TaskType
}

// isPluginEnabled 判断指定插件是否出现在当前配置的启用列表中
func isPluginEnabled(plugin string) bool {

	cfg := config.GetGlobalMisakaConfigure()
	// 当配置启动列表为空的时候 直接返回false
	if cfg == nil || len(cfg.Private.Storage.Plugins) == 0 {
		return false
	}

	plugin = strings.TrimSpace(plugin)
	for _, enabledPlugin := range cfg.Private.Storage.Plugins {
		if strings.TrimSpace(enabledPlugin) == plugin {
			return true // 判断插件是否在启动配置列表中
		}
	}
	return false
}

// normalizePluginName 对插件名做基础归一化处理，删除前后的空格
func normalizePluginName(plugin string) string {
	return strings.TrimSpace(plugin)
}

// markPluginLoaded 标记插件已经完成注册加载。
func markPluginLoaded(plugin string) {
	plugin = normalizePluginName(plugin)
	if plugin == "" {
		return
	}
	loadedPlugins.Store(plugin, true)
}

// GetLoadedPluginsSnapshot 返回当前已经成功注册的插件名称快照。
func GetLoadedPluginsSnapshot() []string {
	plugins := make([]string, 0)
	loadedPlugins.Range(func(key, value any) bool {
		pluginName, ok := key.(string)
		if ok && pluginName != "" {
			plugins = append(plugins, pluginName)
		}
		return true
	})
	sort.Strings(plugins)
	return plugins
}

// RegisterPluginTaskAlias 为指定插件注册业务能力别名到 TaskType 的映射。
func RegisterPluginTaskAlias(plugin string, alias string, taskType tasktype.TaskType) error {
	if !isPluginEnabled(plugin) {
		return nil
	}

        aliasDoc, err := parseTaskAliasDoc(alias, plugin, taskType)
        if err != nil {
                return err
        }
        if aliasDoc.Alias == "" {
		return nil
	}

        taskTypeAlias.Store(aliasDoc.Alias, taskType)
        taskAliasDocs.Store(aliasDoc.Alias, aliasDoc)
	markPluginLoaded(plugin)
	return nil
}

// ResolveTaskType 根据业务能力别名解析实际注册的 TaskType。
func ResolveTaskType(alias string) (tasktype.TaskType, bool) {
        alias = normalizeAliasName(alias)
	if alias == "" {
		return "", false
	}

	value, ok := taskTypeAlias.Load(alias)
	if !ok {
		return "", false
	}

	resolvedTaskType, ok := value.(tasktype.TaskType)
	if !ok {
		return "", false
	}
	return resolvedTaskType, true
}

// GetTaskAliasDocsSnapshot 返回当前已注册的别名规则文档快照。
func GetTaskAliasDocsSnapshot() []TaskAliasDoc {
        docs := make([]TaskAliasDoc, 0)
        taskAliasDocs.Range(func(key, value any) bool {
                doc, ok := value.(TaskAliasDoc)
                if ok {
                        docs = append(docs, doc)
                }
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

func parseTaskAliasDoc(rawAlias string, plugin string, taskType tasktype.TaskType) (TaskAliasDoc, error) {
        rawAlias = strings.TrimSpace(rawAlias)
        if rawAlias == "" {
                return TaskAliasDoc{}, nil
        }

        parts := strings.SplitN(rawAlias, "@", 2)
        alias := normalizeAliasName(parts[0])
        if alias == "" {
                return TaskAliasDoc{}, fmt.Errorf("task alias is empty")
        }

        desc := ""
        if len(parts) == 2 {
                desc = strings.TrimSpace(parts[1])
        }

        return TaskAliasDoc{
                Plugin:   normalizePluginName(plugin),
                Alias:    alias,
                Desc:     desc,
                TaskType: taskType,
        }, nil
}

func normalizeAliasName(alias string) string {
        alias = strings.TrimSpace(alias)
        if alias == "" {
                return ""
        }
        parts := strings.SplitN(alias, "@", 2)
        return strings.TrimSpace(parts[0])
}

// RegisterPluginsActionsDispatcher 根据模块类型分发插件注册动作。
func RegisterPluginsActionsDispatcher(plugins string, module string, action any) error {
	if !isPluginEnabled(plugins) {
		return nil
	}
	markPluginLoaded(plugins)

	switch module {
	case "task":
		return RegisterPluginsActionsInTaskType(plugins, action.(pluginxInterface.FuncTaskTypeAdd))
	default:
		break
	}

	return nil
}

// RegisterPluginsActionsInTaskType 注册插件扩展的 TaskType 列表。
func RegisterPluginsActionsInTaskType(plugins string, action pluginxInterface.FuncTaskTypeAdd) error {
	if !isPluginEnabled(plugins) {
		return nil
	}
	markPluginLoaded(plugins)
	var err error
	if pluginsx.GetPluginsBus().TaskTypes, err = action(pluginsx.GetPluginsBus().TaskTypes); err != nil {
		return err
	}
	return nil
}

// RegisterPluginsActionsInTaskTypeAction 注册指定 TaskType 的执行函数。
func RegisterPluginsActionsInTaskTypeAction(plugins string, taskType tasktype.TaskType, action pluginxInterface.FuncTaskTypeAction) error {
	if !isPluginEnabled(plugins) {
		return nil
	}
	markPluginLoaded(plugins)
	pluginsx.GetPluginsBus().TaskTypeActions.Store(taskType, action)
	return nil
}

// RegisterPluginsActionsInTaskTypeRoll 注册指定 TaskType 的回滚函数。
func RegisterPluginsActionsInTaskTypeRoll(plugin string, taskType tasktype.TaskType, action pluginxInterface.FuncTaskTypeRoll) error {
	if !isPluginEnabled(plugin) {
		return nil
	}
	markPluginLoaded(plugin)
	pluginsx.GetPluginsBus().TaskTypeRoll.Store(taskType, action)
	return nil
}
