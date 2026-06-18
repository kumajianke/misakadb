package atomic_work_center

import (
	"crypto/md5"
	"math/rand"
	"strconv"
	"time"

	"github.com/puzpuzpuz/xsync/v3"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

// 添加作业
func (this *AtomicWorkCenter) AddTask(task *Task, retry int) bool {
	now := time.Now().Format("2006-1-2 15:04:05") + strconv.Itoa(rand.Intn(30000)+1)
	md5Hash := md5.Sum([]byte(now))
	_, loaded := this.TasksMap.LoadOrStore(string(md5Hash[:]), task)
	if loaded {
		// 如果Map中存在
		if retry <= 0 {
			return false
		}
		return this.AddTask(task, retry-1)
	}
	return true
}

// 获取指定作业的ID
func (this *AtomicWorkCenter) GetTask(taskId string) *Task {

	if task, ok := this.TasksMap.Load(taskId); ok {
		return task
	}
	return nil
}

// 让作业继续下一步
func (this *AtomicWorkCenter) DoNext() {

}
