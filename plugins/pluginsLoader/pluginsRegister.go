package pluginsloader

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"misakadb/config"
	pluginsX "misakadb/plugins/pluginsx"
	pluginsXInterface "misakadb/plugins/pluginsx/pluginsx_interface"
	pluginsxInterface "misakadb/plugins/pluginsx/pluginsx_interface"
)

var (
	loadedPlugins sync.Map
)

// isPluginEnabled 判断指定插件是否出现在当前配置的启用列表中。
func isPluginEnabled(plugin string) bool {
	cfg := config.GetGlobalMisakaConfigure()
	if cfg == nil || len(cfg.Private.Storage.Plugins) == 0 {
		return false
	}

	plugin = normalizePluginName(plugin)
	for _, enabledPlugin := range cfg.Private.Storage.Plugins {
		if normalizePluginName(enabledPlugin) == plugin {
			return true
		}
	}
	return false
}

// normalizePluginName 对插件名做基础归一化处理，删除前后的空格。
func normalizePluginName(plugin string) string {
	return NormalizeConfiguredPluginName(plugin)
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
	loadedPlugins.Range(func(key, _ any) bool {
		pluginName, ok := key.(string)
		if ok && pluginName != "" {
			plugins = append(plugins, pluginName)
		}
		return true
	})
	sort.Strings(plugins)
	return plugins
}

// RegisterPluginTaskTypeWithAlias 为指定插件注册业务能力别名到 TaskType 的映射。
func RegisterPluginTaskTypeWithAlias(plugin string, alias string, taskType tasktype.TaskType) error {
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

	markPluginLoaded(plugin)

	PluginsX := pluginsX.GetPluginsX()
	PluginsX.StoreTaskTypeAlias(aliasDoc.Alias, taskType)
	PluginsX.StoreTaskAliasDoc(aliasDoc)
	PluginsX.TaskTypes = append(PluginsX.TaskTypes, taskType)

	return nil
}

// ResolveTaskType 根据业务能力别名解析实际注册的 TaskType。
func ResolveTaskType(alias string) (tasktype.TaskType, bool) {
	alias = normalizeAliasName(alias)
	if alias == "" {
		return "", false
	}
	return pluginsX.GetPluginsX().ResolveTaskType(alias)
}

func parseTaskAliasDoc(rawAlias string, plugin string, taskType tasktype.TaskType) (pluginsX.TaskAliasDoc, error) {
	rawAlias = strings.TrimSpace(rawAlias)
	if rawAlias == "" {
		return pluginsX.TaskAliasDoc{}, nil
	}

	parts := strings.SplitN(rawAlias, "@", 2)
	alias := normalizeAliasName(parts[0])
	if alias == "" {
		return pluginsX.TaskAliasDoc{}, fmt.Errorf("task alias is empty")
	}

	desc := ""
	if len(parts) == 2 {
		desc = strings.TrimSpace(parts[1])
	}

	return pluginsX.TaskAliasDoc{
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
func RegisterPluginsActionsDispatcher(plugin string, module string, action any) error {
	if !isPluginEnabled(plugin) {
		return nil
	}
	markPluginLoaded(plugin)

	switch module {
	case "task":
		return RegisterPluginsActionsInTaskType(plugin, action.(pluginsXInterface.FuncTaskTypeAdd))
	default:
		return nil
	}
}

// RegisterPluginsActionsInTaskType 注册插件扩展的 TaskType 列表。
func RegisterPluginsActionsInTaskType(plugin string, action pluginsXInterface.FuncTaskTypeAdd) error {
	if !isPluginEnabled(plugin) {
		return nil
	}
	markPluginLoaded(plugin)

	PluginsX := pluginsX.GetPluginsX()
	taskTypes, err := action(PluginsX.TaskTypes)
	if err != nil {
		return err
	}
	PluginsX.TaskTypes = taskTypes
	return nil
}

// RegisterPluginsActionsInTaskTypeAction 注册指定 TaskType 的执行函数。
func RegisterPluginsActionsInTaskTypeAction(plugin string, taskType tasktype.TaskType, action pluginsXInterface.FuncTaskTypeAction) error {
	if !isPluginEnabled(plugin) {
		return nil
	}
	markPluginLoaded(plugin)
	pluginsX.GetPluginsX().TaskTypeActions.Store(taskType, action)
	return nil
}

// RegisterPluginsActionsInTaskTypeRoll 注册指定 TaskType 的回滚函数。
func RegisterPluginsActionsInTaskTypeRoll(plugin string, taskType tasktype.TaskType, action pluginsXInterface.FuncTaskTypeRoll) error {
	if !isPluginEnabled(plugin) {
		return nil
	}
	markPluginLoaded(plugin)
	pluginsX.GetPluginsX().TaskTypeRoll.Store(taskType, action)
	return nil
}

func RegisterPluginsTaskCombo(combo_name string, combo_func pluginsxInterface.FuncTaskCombo) error {
	pluginsx := pluginsX.GetPluginsX()
	if _, ok := pluginsx.TaskCombo.Load(combo_name); ok {
		return fmt.Errorf("has duplicate combo name: %s", combo_name)
	}
	pluginsx.TaskCombo.Store(combo_name, combo_func)
	return nil
}

func GetPluginsTaskCombo(combo_name string) (pluginsXInterface.FuncTaskCombo, error) {
	pluginx := pluginsX.GetPluginsX()
	combo, ok := pluginx.TaskCombo.Load(combo_name)
	if !ok {
		return nil, errors.New("No such as the key in pluginx.TaskCombo!")
	}
	return combo, nil
}
