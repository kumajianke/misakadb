package atomic_work_center_test

import (
	atomic_work_center "misakadb/atomic/atomicWorkCenter"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
)

func TestTask() {
	atomicWorkCenter := atomic_work_center.NewAtomicWorkCenter()
	task := atomic_work_center.NewTask(nil)

	task.TaskBooks = tasktype.NewShipBuilder().Add(
		tasktype.TaskModFile,
		"xx.file",
	).Build()
	ok, task_id := atomicWorkCenter.AddTask(task, 3)
	if ok {
		atomicWorkCenter.DoNext(task_id)
	}
}
