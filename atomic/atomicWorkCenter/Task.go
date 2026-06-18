package atomic_work_center

import (
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"time"
)

type EnumTaskStatus int

const (
	Pending  EnumTaskStatus = iota // 待处理
	Running                        // 运行中
	Success                        // 完成
	BackRoll                       // 回滚
)

type Task struct {
	TaskStatus      EnumTaskStatus
	TaskReleaseTime time.Time
	TaskBody        []*tasktype.TaskTypeShip
}

func NewTask(taskBody []*tasktype.TaskTypeShip) *Task {
	return &Task{
		TaskStatus:      Pending,
		TaskReleaseTime: time.Now().Add(time.Second * 10), // 十秒后过期
		TaskBody:        taskBody,
	}
}
