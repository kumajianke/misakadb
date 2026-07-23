package pluginsloader

import (
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	pluginsx "misakadb/plugins/pluginx"
)

// 插件注册分发器
func RegisterPluginsActionsDispatcher(plugins string, module string, action any) error {

	switch module {
	case "task":
		return RegisterPluginsActionsInTaskType(plugins, action.(func(allTaskTypelst []tasktype.TaskType) ([]tasktype.TaskType, error)))
	default:
		break
	}

	return nil
}

// 用于注册TaskType(传入方法句柄)
func RegisterPluginsActionsInTaskType(plugins string, action func(allTaskTypelst []tasktype.TaskType) ([]tasktype.TaskType, error)) error {
	var err error

	if pluginsx.GetPluginsBus().AllTaskType, err = action(pluginsx.GetPluginsBus().AllTaskType); err != nil {
		return err
	}

	return nil
}
