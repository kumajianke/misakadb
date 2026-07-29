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

func OnRemoveFile(taskType tasktype.TaskType, params []string) error {
	// 实现对文件的删除操作
	return nil
}

func RollRemoveFile(taskType tasktype.TaskType, params []string) error {
	// 回滚删除操作
	return nil
}

/*
添加TaskType
*/
func AddTaskType() error {
	if err := pluginsloader.RegisterPluginTaskTypeWithAlias(modName, "misaka.removefile@用于删除文件的tasktype", TaskRemoveFile); err != nil {
		return fmt.Errorf("register alias %s failed: %w", "misaka.removefile", err)
	}
	if err := pluginsloader.RegisterPluginTaskTypeWithAlias(modName, "misaka.removefolder@用于删目录的tasktype", TaskRemoveFolder); err != nil {
		return fmt.Errorf("register alias %s failed: %w", "misaka.removefolder", err)
	}
	return nil
}

/*
添加TaskTypeAction
*/
func AddTaskTypeAction() error {
	if err := pluginsloader.RegisterPluginsActionsInTaskTypeAction(modName, TaskRemoveFile, OnRemoveFile); err != nil {
		return fmt.Errorf("register task action failed: %w", err)
	}
	return nil
}

func Register() error {
	if err := AddTaskType(); err != nil {
		return err
	}

	if err := AddTaskTypeAction(); err != nil {
		return err
	}

	clilog.Success("基础插件加载完毕.")
	return nil
}
