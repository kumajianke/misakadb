package eventbus

import "sync"

/**
 * 原子任务中心事件总线 用于原子任务中心的内部通讯
 */

type AtomicWorkEventBus struct {
	EventBus chan string
}

var AtomicWorkEventBusInstance *AtomicWorkEventBus
var atomicWorkEventBusOnce sync.Once

func GetAtomicWorkCenterEventBus() *AtomicWorkEventBus {
	atomicWorkEventBusOnce.Do(func() {
		AtomicWorkEventBusInstance = &AtomicWorkEventBus{
			EventBus: make(chan string, 10),
		}
	})
	return AtomicWorkEventBusInstance
}
