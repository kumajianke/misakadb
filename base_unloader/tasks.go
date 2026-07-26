package base_unloader

import (
	"fmt"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"misakadb/clilog"
	pluginsloader "misakadb/plugins/pluginsLoader"
)

const (
	TaskRemoveFile   tasktype.TaskType = "remove_file"
	TaskModFile      tasktype.TaskType = "mod_file"
	TaskRemoveFolder tasktype.TaskType = "remove_folder"
)

func AddTaskType(allTaskTypelst []tasktype.TaskType) ([]tasktype.TaskType, error) {

	allTaskTypelst = append(allTaskTypelst, TaskRemoveFile)
	allTaskTypelst = append(allTaskTypelst, TaskModFile)
	allTaskTypelst = append(allTaskTypelst, TaskRemoveFolder)
	return allTaskTypelst, nil
}

func OnRemoveFile(taskType tasktype.TaskType, params []string) error {
	// 实现对文件的删除操作
	return nil

}

func RollRemoveFile(taskType tasktype.TaskType, params []string) error {
	// 回滚删除操作
	return nil
}

func init() {
	pluginsloader.RegisterBuiltinPlugin(modName, Register)
}

func Register() error {
	if err := pluginsloader.RegisterPluginTaskAlias(modName, "misaka.removefile@用于删除文件的tasktype", TaskRemoveFile); err != nil {
		return fmt.Errorf("register alias %s failed: %w", "misaka.removefile", err)
	}
	if err := pluginsloader.RegisterPluginTaskAlias(modName, "misaka.removefolder@用于删目录的tasktype", TaskRemoveFolder); err != nil {
		return fmt.Errorf("register alias %s failed: %w", "misaka.removefolder", err)
	}
	if err := pluginsloader.RegisterPluginsActionsInTaskType(modName, AddTaskType); err != nil {
		return fmt.Errorf("register task types failed: %w", err)
	}

	if err := pluginsloader.RegisterPluginsActionsInTaskTypeAction(modName, TaskRemoveFile, OnRemoveFile); err != nil {
		return fmt.Errorf("register task action failed: %w", err)
	}
	clilog.Success("基础插件加载完毕.")
	return nil
}
