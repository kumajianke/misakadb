package atomic_work_center_test

import (
	atomic_work_center "misakadb/atomic/atomicWorkCenter"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
)

func TestTask() {
	atomicWorkCenter := atomic_work_center.NewAtomicWorkCenter()
	task := atomic_work_center.NewTask(nil)
	task.TaskBody = []*tasktype.TaskTypeShip{
		tasktype.NewTaskTypeShip(
			tasktype.TaskRemoveFile,
			"test.txt",
		),
	}
	atomicWorkCenter.AddTask(task, 3)
}
