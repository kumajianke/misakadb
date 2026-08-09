package atomic_work_center

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"
	eventbus "misakadb/atomic/EventBus"
	tasktype "misakadb/atomic/atomicWorkCenter/TaskType"
	"misakadb/clilog"
	"misakadb/lock/global_lock"
	pluginsx "misakadb/plugins/pluginsX"
	"strconv"
	"time"

	"github.com/puzpuzpuz/xsync/v3"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 作业中心
type AtomicWorkCenter struct {
	TasksMap *xsync.MapOf[string, *Task]
}

var AtomicWorkCenterInstance *AtomicWorkCenter

/*
DESCRIPTION

	新建原子任务中心，这个是单例存在的，并且线程安全

PARAMS

	NULL
*/
func NewAtomicWorkCenter() *AtomicWorkCenter {
	if AtomicWorkCenterInstance == nil {
		AtomicWorkCenterInstance = &AtomicWorkCenter{
			TasksMap: xsync.NewMapOf[string, *Task](),
		}
	}
	return AtomicWorkCenterInstance
}

/*
DESCRIPTION

	向原子任务中心添加任务

Params

	task: 任务体信息
	retry: 重试次数，一般填写为3即可
*/
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

/*
DESCRIPTION

	获取指定作业的ID

Parmas

	0 任务的ID，通过AddTask获取的加密ID
*/
func (this *AtomicWorkCenter) GetTask(taskId string) *Task {

	if task, ok := this.TasksMap.Load(taskId); ok {
		return task
	}
	return nil
}

/*
DESCRIPTION

	获取当前Task的Index指向的TaskBook

PARAMS

	0 taskID  任务的ID 通过AddTask获取加密的ID
*/
func (this *AtomicWorkCenter) GetCurrentTaskBookRef(taskId string) (*tasktype.TaskBooks, error) {
	task := this.GetTask(taskId) // 获取任务信息
	index := task.TaskCurrentIndex
	if task.TaskCurrentIndex > (len(task.TaskBooks) - 1) {
		return tasktype.NewTaskBooks("invalid TaskBooks"), errors.New("Index range out the length of TaskBook")
	}
	current_task := task.TaskBooks[index]
	return current_task, nil

}

/*
DESCRIPTION

	让作业继续下一步，并对异常状态的作业进行处理

RETURNS

	error 报错，当执行Action报错的时候会标记为回滚

PARAMS

	taskID: 任务ID 通过AddTask获取
*/
func (this *AtomicWorkCenter) DoNext(taskId string) error {
	task := this.GetTask(taskId)
	fmt.Println(task)

	if task == nil {
		return errors.New("[AtomicWorkCenter] KeyError No such an the keyvalue in map!")
	}

	// 长度验证
	if task.TaskCurrentIndex+1 >= len(task.TaskBooks) {
		return errors.New("[AtomicWorkCenter] Index Error!")
	}

	// 作业状态信息验证
	if task.TaskStatus == Pending && task.TaskCurrentIndex == 0 {
		// 状态为等待状态修改为Running
		task.TaskStatus = Running
	} else if task.TaskStatus == BackRoll {
		// TODO: 回滚机制
		this.CancleTask(taskId) // 取消代码
	} else if task.TaskStatus == Success {
		// 已完成的任务本

	}

	if task.TaskStatus == Running {
		// 运行下一个函数

		pluginx_bus := pluginsx.GetPluginsX()
		taskbooks, err := this.GetCurrentTaskBookRef(taskId)
		taskbooks_tasktype := taskbooks.TaskType

		if err != nil {
			return err
		}

		// TaskBooks 存储了tasktype的别名 通过别名获取对应的tasktype
		tasktype, ok := pluginx_bus.TaskTypeAlias.Load(string(taskbooks_tasktype)) // 获取TaskType对应的真实TaskType
		if !ok {
			return errors.New("can not get the tasktype from pluginx by tasktype of taskbooks.")
		}

		// 通过tasktype在插件中获取对应的tasktype的执行方案(tasktype actions)
		actions_fn, actions_ok := pluginx_bus.TaskTypeActions.Load(tasktype)
		if !actions_ok {
			return errors.New("can not to load the actions of tasktype which named " + string(tasktype))
		}

		// 执行 tasktype actions
		err = actions_fn(tasktype, taskbooks.Params)
		if err != nil {
			task.TaskStatus = BackRoll // 待回滚
			return err                 // action函数本身的报错， 这里报错的化需要回滚整个任务链
		}
	}
	var eb *eventbus.AtomicWorkEventBus
	eb = eventbus.NewAtomicWorkCenterEventBus()
	eb.EventBus <- "sync-to-local" // 同步任务信息
	return nil
}

/*
DESCRIPTIONS

	（不安全）取消某个任务，线程安全，但不进行回滚

PARAMS

	ID 任务ID 通过AddTask获取
*/
func (this *AtomicWorkCenter) NoSafeRemoveTask(taskId string) *Task {
	task, _ := this.TasksMap.LoadAndDelete(taskId)
	return task
}

// TODO: 取消任务，删除并回滚
func (this *AtomicWorkCenter) CancleTask(taskId string) *Task {
	return nil
}

// TODO: 完成所有的任务的作业本 直到完成
/*
DESCRIPTION
	完成指定的任务计划（单线程）并返回是否成功
PARAMS
	当前的任务ID，通过AddTask的时候获取
RETURNS
	error: 函数执行的时候出现的错误
	bool : 当前任务是否被成功执行到结束
*/
func (this *AtomicWorkCenter) DoSustain(taskId string) (error, bool) {
	task := this.GetTask(taskId)

	// 任务合规性验证
	if task.TaskReleaseTime.After(time.Now()) { // 任务过期的情况
		this.NoSafeRemoveTask(taskId)
		return errors.New("the task is release now."), false
	}
	if task.TaskStatus == Success { // 任务已经被执行成功了
		this.NoSafeRemoveTask(taskId)
		return errors.New("the task is finish."), true
	}

	// 任务处理开始
	if task == nil {
		// 任务不存在的时候额处理方式
		return errors.New("[AtomicWorkCenter] KeyError No such an the keyvalue in map!"), false
	}
	_, unlock_write_lock_atomic_work_center, err := global_lock.GetOrStoreGlobalLock(
		"write_lock_atomic_work_center", "lock",
	) // 对当前的原子任务中心进行加锁 防止线程冲突

	if err != nil {
		return errors.New("[AtomicWorkCenter] acquire write lock failed：" + err.Error()), false
	} // 独占写锁
	defer unlock_write_lock_atomic_work_center() // 解写锁

	for task.TaskCurrentIndex < len(task.TaskBooks) {
		err := this.DoNext(taskId) // 执行当前任务的下一个作业
		if err != nil {
			clilog.Error(err)
			this.CancleTask(taskId) // 回滚任务
		}
		task.TaskCurrentIndex++                      // 索引增加以便执行下一个作业
		eb := eventbus.NewAtomicWorkCenterEventBus() // 创建新的事件总线
		eb.EventBus <- "sync-to-local"               // 序列化内容到服务器本地 以防止宕机恢复
	}

	return nil, true
}
