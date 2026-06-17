package workcenterserializer

import (
	"encoding/json"
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
	_, unlock, err := global_lock.GetOrStoreGlobalLock("WorkCenterSerializerDump", "l")
	if err != nil {
		this._thread__chan <- "file_write_error"
		return
	}
	defer unlock()

	if _, err = os.Stat(".data"); err != nil {
		err = os.Mkdir(".data", 0700)
		if err != nil {
			this._thread__chan <- "file_write_error"
			return
		}
	}

	jsonData, err := json.Marshal(this)
	if err != nil {
		this._thread__chan <- "file_write_error"
		return
	}
	err = os.WriteFile(".data/work_center.json", jsonData, 0700)
	if err != nil {
		this._thread__chan <- "file_write_error"
		return
	}
	this._thread__chan <- "dump_success"
}

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

func FastInitWorkCenterSerializer() *WorkCenterSerializer {
	workCenterSerializer := BuildWorkCenterSerializer(
		atomic_work_center.NewAtomicWorkCenter(),
	)
	workCenterSerializer.WorkerCenterSerializerThread(3)
	return workCenterSerializer
}
