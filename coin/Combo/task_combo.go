package base_combo

import (
	"errors"
	atomic_work_center "misakadb/atomic/atomicWorkCenter"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	baseconfig "misakadb/coin/base_config"
	"misakadb/config"
	"misakadb/lock/global_lock"
	"path"
)

/*
DESCRIPTION

	删除数据库的TaskCombo

PARAMS

	0 需要删除的数据库库名

RETURNS

	error 报错信息
	string Combo对应的TaskID
*/
func ComboRemoveDatabase(params []string) (error, string) {
	if len(params) != 1 {
		return errors.New("Argument Error on ComboRemoveDatabase!"), ""
	}

	awc := atomic_work_center.GetAtomicWorkCenter()
	storage := config.GetGlobalMisakaConfigure().Private.Storage
	db_name := params[0]
	db_path := path.Join(storage.Path, db_name)
	_, unlock, err := global_lock.GetOrStoreGlobalLock(
		global_lock.GetLocksPrefix().DBFile+db_name,
		"l",
	)
	if err != nil {
		return err, ""
	}
	defer unlock()

	taskbook := tasktype.NewShipBuilder().Add(baseconfig.TaskRemoveFolder, db_path).Build()
	task := atomic_work_center.NewTask(taskbook)
	ok, task_id := awc.AddTask(task, 3)
	if !ok {
		return errors.New("create task failed"), ""
	}

	// 执行删除任务书
	if err, ok := awc.DoSustain(task_id); err != nil || !ok {
		if err == nil {
			err = errors.New("execute delete database task failed")
		}
		return err, task_id
	}
	return nil, task_id
}
