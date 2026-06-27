package atomic_work_center

import (
	"crypto/md5"
	"math/rand"
	eventbus "misakadb/atomic/EventBus"
	"misakadb/clilog"
	"misakadb/lock/global_lock"
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
func (this *AtomicWorkCenter) AddTask(task *Task, retry int) (bool, string) {
	// 加上写锁
	_, unlock_write_lock_atomic_work_center, err := global_lock.GetOrStoreGlobalLock("write_lock_atomic_work_center", "lock")
	if err != nil {
		clilog.Error("[AtomicWorkCenter.AddTask] acquire write lock failed:", err)
		return false, ""
	} // 独占写锁
	defer unlock_write_lock_atomic_work_center()

	now := time.Now().Format("2006-1-2 15:04:05") + strconv.Itoa(rand.Intn(30000)+1)
	md5Hash := md5.Sum([]byte(now))
	string_md5 := string(md5Hash[:])
	_, loaded := this.TasksMap.LoadOrStore(string_md5, task)
	if loaded {
		// 如果Map中存在
		if retry <= 0 {
			return false, ""
		}
		return this.AddTask(task, retry-1)
	}
	return true, string_md5
}

// 获取指定作业的ID
func (this *AtomicWorkCenter) GetTask(taskId string) *Task {

	if task, ok := this.TasksMap.Load(taskId); ok {
		return task
	}
	return nil
}

// 让作业继续下一步
func (this *AtomicWorkCenter) DoNext(taskId string) {
	task := this.GetTask(taskId)

	if task == nil {
		clilog.Error("[AtomicWorkCenter] KeyError No such an the keyvalue in map!")
		return
	}

	// 长度验证
	if task.TaskCurrentIndex+1 >= len(task.TaskBody) {
		clilog.Error("[AtomicWorkCenter] Index Error!")
	}

	// 作业状态信息验证
	if task.TaskStatus == Pending && task.TaskCurrentIndex == 0 {
		// 状态为等待状态修改为Running
		task.TaskStatus = Running
	} else if task.TaskStatus == BackRoll {
		// TODO 回滚机制
	} else if task.TaskStatus == Success {
		this.TasksMap.LoadAndDelete(taskId) // 任务完成 需要进行删除操作
	} else if task.TaskStatus == Running {
		// TODO 运行下一个函数
	}
	var eb *eventbus.AtomicWorkEventBus
	eb = eventbus.NewAtomicWorkCenterEventBus()
	eb.EventBus <- "sync-to-local"
}

// TODO 持续交付一个作业 直到完成
func (this *AtomicWorkCenter) DoSustain(taskId string) {

}
