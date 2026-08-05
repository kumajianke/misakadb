package atomic_work_center

import (
	"crypto/md5"
	"errors"
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
	_, unlock_write_lock_atomic_work_center, err := global_lock.GetOrStoreGlobalLock(
		"write_lock_atomic_work_center", "lock",
	)
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
	if task.TaskCurrentIndex+1 >= len(task.TaskBooks) {
		clilog.Error("[AtomicWorkCenter] Index Error!")
	}

	// 作业状态信息验证
	if task.TaskStatus == Pending && task.TaskCurrentIndex == 0 {
		// 状态为等待状态修改为Running
		task.TaskStatus = Running
	} else if task.TaskStatus == BackRoll {
		// TODO: 回滚机制
	} else if task.TaskStatus == Success {

	}

	if task.TaskStatus == Running {
		// TODO: 运行下一个函数
	}
	var eb *eventbus.AtomicWorkEventBus
	eb = eventbus.NewAtomicWorkCenterEventBus()
	eb.EventBus <- "sync-to-local"
}

func (this *AtomicWorkCenter) RemoveTask(taskId string) *Task {
	task, _ := this.TasksMap.LoadAndDelete(taskId)
	return task
}

// TODO: 完成所有的任务的作业本 直到完成
func (this *AtomicWorkCenter) DoSustain(taskId string) error {
	task := this.GetTask(taskId)

	// 任务合规性验证
	if task.TaskReleaseTime.After(time.Now()) {
		this.RemoveTask(taskId)
		return errors.New("the task is release now.")
	}
	if task.TaskStatus == Success {
		this.RemoveTask(taskId)
		return errors.New("the task is finish.")
	}

	// 任务处理开始
	if task == nil {
		return errors.New("[AtomicWorkCenter] KeyError No such an the keyvalue in map!")
	}
	_, unlock_write_lock_atomic_work_center, err := global_lock.GetOrStoreGlobalLock(
		"write_lock_atomic_work_center", "lock",
	)
	if err != nil {
		return errors.New("[AtomicWorkCenter] acquire write lock failed：" + err.Error())
	} // 独占写锁
	defer unlock_write_lock_atomic_work_center() // 解写锁

	for task.TaskCurrentIndex < len(task.TaskBooks) {
		this.DoNext(taskId)                          // 执行当前任务的下一个作业
		task.TaskCurrentIndex++                      // 索引增加以便执行下一个作业
		eb := eventbus.NewAtomicWorkCenterEventBus() // 创建新的事件总线
		eb.EventBus <- "sync-to-local"               // 序列化内容到服务器本地 以防止宕机恢复
	}

	return nil
}
