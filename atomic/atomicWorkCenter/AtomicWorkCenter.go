package atomic_work_center

import (
	"github.com/puzpuzpuz/xsync/v3"
)

// 作业中心
type AtomicWorkCenter struct {
	TasksMap *xsync.MapOf[string, *Task]
}
