package atomic_work_center_test

import (
	atomic_work_center "misakadb/atomic/atomicWorkCenter"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"misakadb/clilog"
	pluginsloader "misakadb/plugins/pluginsLoader"
)

func TestTask() {
	atomicWorkCenter := atomic_work_center.NewAtomicWorkCenter()
	task := atomic_work_center.NewTask(nil)

	remove_file, ok := pluginsloader.ResolveTaskType(
		"misaka.removefile",
	)

	if !ok {
		clilog.Error("[ERROR] can not get the tasktype!")
	}
	task.TaskBooks = tasktype.NewShipBuilder().Add(
		remove_file,
		"xx.file",
	).Build()
	ok, task_id := atomicWorkCenter.AddTask(task, 3)
	if ok {
		atomicWorkCenter.DoNext(task_id)
	}
}
