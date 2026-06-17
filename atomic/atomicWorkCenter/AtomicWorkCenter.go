package atomic_work_center

import (
	"github.com/puzpuzpuz/xsync/v3"
)

var AtomicWorkCenterInstance *AtomicWorkCenter

// 作业中心
type AtomicWorkCenter struct {
	TasksMap *xsync.MapOf[string, *Task]
}

func NewAtomicWorkCenter() *AtomicWorkCenter {
	if AtomicWorkCenterInstance == nil {
		AtomicWorkCenterInstance = &AtomicWorkCenter{
			TasksMap: xsync.NewMapOf[string, *Task](),
		}
	}
	return AtomicWorkCenterInstance
}
