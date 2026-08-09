package base_unloader

import (
	"errors"
	"fmt"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	baseconfig "misakadb/base_unloader/base_config"
	"misakadb/lock/global_lock"
	pluginsloader "misakadb/plugins/pluginsLoader"
	"misakadb/plugins/pluginsx"
	pluginsxInterface "misakadb/plugins/pluginsx/pluginsx_interface"
	filesafe "misakadb/shares/filesafe"
	"os"
	"path/filepath"
)

/*
DESCRIPTION

	实现对文件的删除操作 用于TaskType事件

Params 0 对应文件的路径
*/
func OnRemoveFile(taskType tasktype.TaskType, params []string) error {
	if len(params) != 1 {
		return errors.New("Argument error: you can research onRemoveFile on code project to get params comments!")
	}

	handle_file_path := params[0]
	only_path := filepath.Dir(handle_file_path)  // 路径
	only_name := filepath.Base(handle_file_path) // 名称
	unlocked, err := global_lock.LockFileHandle(handle_file_path)
	if err != nil {
		return err
	}
	defer unlocked()

	remove_state := os.Rename(handle_file_path, filepath.Join(only_path, "MisakaRemove."+only_name)) // 尝试修改
	if remove_state != nil {
		return remove_state
	}
	return filesafe.Fsync(only_path) // 强制落盘
}

/*
DESCRIPTION

	TaskType的删除文件的回滚操作

Params 0 对应的文件路径（不需要MisakaRemove.前缀)
*/
func RollRemoveFile(taskType tasktype.TaskType, params []string) error {
	// 回滚删除操作
	handle_file_path := params[0]
	only_path := filepath.Dir(handle_file_path)
	only_name := filepath.Base(handle_file_path)
	remove_name := "MisakaRemove." + only_name

	remove_path := filepath.Join(only_path, remove_name)
	unlocked, err := global_lock.LockFileHandle(remove_path)
	if err != nil {
		return err
	}
	defer unlocked()
	undo_state := os.Rename(remove_path, handle_file_path)

	if undo_state != nil {
		return undo_state
	}

	return filesafe.Fsync(only_path)
}

/*
添加TaskType
*/
func AddTaskType() error {
	if err := pluginsloader.RegisterPluginTaskTypeWithAlias(baseconfig.ModName, "misaka.removefile@用于删除文件的tasktype", baseconfig.TaskRemoveFile); err != nil {
		return fmt.Errorf("register alias %s failed: %w", "misaka.removefile", err)
	}
	if err := pluginsloader.RegisterPluginTaskTypeWithAlias(baseconfig.ModName, "misaka.removefolder@用于删目录的tasktype", baseconfig.TaskRemoveFolder); err != nil {
		return fmt.Errorf("register alias %s failed: %w", "misaka.removefolder", err)
	}
	return nil
}

/*
添加TaskTypeAction
*/
func AddTaskTypeAction() error {
	if err := pluginsloader.RegisterPluginsActionsInTaskTypeAction(baseconfig.ModName, baseconfig.TaskRemoveFile, OnRemoveFile); err != nil {
		return fmt.Errorf("register task action failed: %w", err)
	}
	return nil
}

/*
添加 TaskCombo
*/
func AddTaskCombo(combo_name string, combo_func pluginsxInterface.FuncTaskCombo) error {
	pluginsx.GetPluginsX()
	return nil
}
