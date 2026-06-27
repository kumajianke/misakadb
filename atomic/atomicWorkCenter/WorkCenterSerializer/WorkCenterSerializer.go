package work_center_serializer

import (
	"encoding/json"
	eventbus "misakadb/atomic/EventBus"
	atomic_file_handler "misakadb/atomic/atomicFileHandler"
	atomic_work_center "misakadb/atomic/atomicWorkCenter"
	"misakadb/clilog"
	"misakadb/lock/global_lock"
	"os"
	"time"
)

type WorkCenterSerializer struct {
	TaskMap       map[string]*atomic_work_center.Task `json:"task_map"`
	_thread__chan chan string                         `json:"-"`
}

func BuildWorkCenterSerializer(wc *atomic_work_center.AtomicWorkCenter) *WorkCenterSerializer {
	buildWorkCenterSerializer := &WorkCenterSerializer{
		TaskMap:       make(map[string]*atomic_work_center.Task),
		_thread__chan: make(chan string, 5),
	} // 创建一个序列化容器

	syncMap := wc.TasksMap // 作业中心任务映射

	// 遍历任务映射，将任务添加到序列化器中
	syncMap.Range(func(key string, value *atomic_work_center.Task) bool {
		buildWorkCenterSerializer.TaskMap[key] = value
		return true
	})
	return buildWorkCenterSerializer
}

func (this *WorkCenterSerializer) Dump() {
	// 加锁
	_, unlock, err := global_lock.GetOrStoreGlobalLock("WorkCenterSerializerDump", "l")
	if err != nil {
		this._thread__chan <- "file_write_error"
		return
	}
	defer unlock()

	// 创建文件夹
	if _, err = os.Stat(".data"); err != nil {
		err = os.Mkdir(".data", 0700)
		if err != nil {
			this._thread__chan <- "file_write_error"
			return
		}
	}

	// 存储Marshal化数据
	jsonData, err := json.Marshal(this)

	_, err = atomic_file_handler.ChunkAtomicSyncWriteFile(
		".data/work_center.json",
		jsonData,
		0644,
	)

	if err != nil {
		this._thread__chan <- "file_write_error"
		return
	}
	// 强制刷盘

	if err != nil {
		this._thread__chan <- "file_write_error"
		return
	}

	this._thread__chan <- "dump_success"
}

// 启动协程进行存储 并通过队列监听是否完成
func (this *WorkCenterSerializer) WorkerCenterSerializerThread(retryTimes int) bool {
	go this.Dump()
	select {
	case msg := <-this._thread__chan:
		if msg == "dump_success" {
			return true
		} else {
			if retryTimes > 0 {
				clilog.Error("[dumpwork-center]" + msg)
				return this.WorkerCenterSerializerThread(retryTimes - 1)
			} else {
				return false
			}
		}
	case <-time.After(100 * time.Second):
		return false
	}
}

// 快速启动原子作业持久化
func FastInitWorkCenterSerializer() *WorkCenterSerializer {
	workCenterSerializer := BuildWorkCenterSerializer(
		atomic_work_center.NewAtomicWorkCenter(),
	)
	eb := eventbus.NewAtomicWorkCenterEventBus()
	go func() {
		clilog.Success("[Atomic Work Center Serializer]fast boot over!")
		for msg := range eb.EventBus {
			if msg == "sync-to-local" {
				workCenterSerializer.WorkerCenterSerializerThread(3)
			}
		}
	}()
	return workCenterSerializer
}

// 加载本地持久化信息
func LoadWorkCenterSerializer() *WorkCenterSerializer {
	workCenterSerializer := &WorkCenterSerializer{}

	jsonData, err := atomic_file_handler.ChunkRead(".data/work_center.json")
	if err != nil {
		clilog.Warning("No File To Load")
		return nil
	}

	if err := json.Unmarshal(jsonData, workCenterSerializer); err != nil {
		clilog.Error("[loadwork-center]" + err.Error())
		return nil
	}

	return workCenterSerializer
}

// 将序列化器转换为原子作业中心
func (this *WorkCenterSerializer) Convert2AtomicWorkCenter() (
	*atomic_work_center.AtomicWorkCenter,
	error,
) {
	clilog.Warning(
		"Convert2AtomicWorkCenter Will Overwrite All Tasks " +
			"In Memory AtomicWorkCenter ! Please Be Careful!",
	)

	// 加上写锁: 1. 序列化器写锁 2. 原子作业中心写锁
	_, unlock_write_lock_atomic_work_center_serializer, err_ser := global_lock.GetOrStoreGlobalLock("write_lock_atomic_work_center_serializer", "lock")
	if err_ser != nil {
		clilog.Error("[Convert2AtomicWorkCenter] acquire serializer lock failed:", err_ser)
		return nil, err_ser
	}
	defer unlock_write_lock_atomic_work_center_serializer()

	_, unlock_write_lock_atomic_work_center, err := global_lock.GetOrStoreGlobalLock("write_lock_atomic_work_center", "lock")

	if err != nil {
		clilog.Error("[Convert2AtomicWorkCenter] acquire atomic work center lock failed:", err)
		return nil, err
	}
	defer unlock_write_lock_atomic_work_center()

	taskMap := this.TaskMap // 序列器所有的map
	atomicWorkCenter := atomic_work_center.NewAtomicWorkCenter()

	atomicWorkCenter.TasksMap.Clear() // 清空原子作业中心任务映射

	for key, task := range taskMap {
		atomicWorkCenter.TasksMap.Store(key, task) // 添加任务到原子作业中心任务映射
	}

	return atomicWorkCenter, nil
}
