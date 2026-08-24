package coin

import (
	"errors"
	"fmt"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"misakadb/clilog"
	base_combo "misakadb/coin/Combo"
	baseconfig "misakadb/coin/base_config"
	"misakadb/lock/global_lock"
	pluginsloader "misakadb/plugins/pluginsLoader"
	"misakadb/shares"
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

func OnRemoveFolder(taskType tasktype.TaskType, params []string) error {
	if len(params) != 1 {
		return errors.New("Argument error: remove folder requires one path")
	}
	path := params[0]
	parent, name := filepath.Dir(path), filepath.Base(path)
	unlock, err := global_lock.LockFileHandle(path)

	if err != nil {
		return err
	}

	defer unlock()
	new_path := filepath.Join(parent, "MisakaRemove."+name)

	f, err := os.Open(new_path)
	if err == nil {
		// 如果文件已经存在那里把她删除掉
		f.Close()
		clilog.Warning("[Warning] The cache generated for the deleted file conflicts with the cache of the historical file, now remove the history")
		os.RemoveAll(new_path)
	}

	if err := os.Rename(path, new_path); err != nil {
		clilog.Error("rename the folder error:" + err.Error())
		return err
	}

	if shares.IsWindows() {
		return nil
	}

	return filesafe.Fsync(parent)
}

func RollRemoveFolder(taskType tasktype.TaskType, params []string) error {
	if len(params) != 1 {
		return errors.New("Argument error: rollback remove folder requires one path")
	}
	path := params[0]
	parent, name := filepath.Dir(path), filepath.Base(path)
	removed := filepath.Join(parent, "MisakaRemove."+name)
	unlock, err := global_lock.LockFileHandle(removed)
	if err != nil {
		return err
	}
	defer unlock()
	if err := os.Rename(removed, path); err != nil {
		return err
	}
	return filesafe.Fsync(parent)
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
	if err := pluginsloader.RegisterPluginsActionsInTaskTypeAction(baseconfig.ModName, baseconfig.TaskRemoveFolder, OnRemoveFolder); err != nil {
		return fmt.Errorf("register folder task action failed: %w", err)
	}
	if err := pluginsloader.RegisterPluginsActionsInTaskTypeRoll(baseconfig.ModName, baseconfig.TaskRemoveFile, RollRemoveFile); err != nil {
		return fmt.Errorf("register task rollback failed: %w", err)
	}
	if err := pluginsloader.RegisterPluginsActionsInTaskTypeRoll(baseconfig.ModName, baseconfig.TaskRemoveFolder, RollRemoveFolder); err != nil {
		return fmt.Errorf("register folder rollback failed: %w", err)
	}
	return nil
}

/*
添加 TaskCombo
*/
func AddTaskCombo() error {

	if err := pluginsloader.RegisterPluginsTaskCombo("drop_db", base_combo.ComboRemoveDatabase); err != nil {
		return err
	} // 添加删除数据库的 Combo

	return nil
}
