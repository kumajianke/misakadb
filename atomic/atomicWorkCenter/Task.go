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
	TaskCurrentIndex int                   //  当前任务执行到了哪里
	TaskStatus       EnumTaskStatus        // 当前任务的状态
	TaskReleaseTime  time.Time             // 当前任务在什么时候回直接释放
	TaskBooks        []*tasktype.TaskBooks // 当前任务的作业本
}

func NewTask(TaskBooks []*tasktype.TaskBooks) *Task {
	return &Task{
		TaskStatus:      Pending,
		TaskReleaseTime: time.Now().Add(time.Second * 10), // 十秒后过期
		TaskBooks:       TaskBooks,
	}
}
