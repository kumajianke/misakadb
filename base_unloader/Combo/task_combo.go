package combo

import (
	"errors"
	atomic_work_center "misakadb/atomic/atomicWorkCenter"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	baseconfig "misakadb/base_unloader/base_config"
	"misakadb/config"
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

	awc := atomic_work_center.NewAtomicWorkCenter()
	storage := config.GetGlobalMisakaConfigure().Private.Storage
	db_name := params[0]
	db_path := path.Join(storage.Path, db_name)
	taskbook := tasktype.NewShipBuilder().Add(baseconfig.TaskRemoveFile, db_path).Build()
	task := atomic_work_center.NewTask(taskbook)
	ok, task_id := awc.AddTask(task, 3)
	if ok {
		return nil, task_id
	}
	return errors.New("Task Add Error on ComboRemoveDatabase!"), ""
}
