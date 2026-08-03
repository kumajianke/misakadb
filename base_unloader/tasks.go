package base_unloader

import (
	"errors"
	"fmt"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"misakadb/clilog"
	"misakadb/lock/global_lock"
	pluginsloader "misakadb/plugins/pluginsLoader"
	filesafe "misakadb/shares/filesafe"
	"os"
	"path/filepath"
)

const (
	TaskRemoveFile   tasktype.TaskType = "remove_file"
	TaskModFile      tasktype.TaskType = "mod_file"
	TaskRemoveFolder tasktype.TaskType = "remove_folder"
)

/*
DESCRIPTION

	实现对文件的删除操作

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
