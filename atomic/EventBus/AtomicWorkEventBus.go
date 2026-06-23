package eventbus

/**
 * 原子任务中心事件总线 用于原子任务中心的内部通讯
 */

type AtomicWorkEventBus struct {
	EventBus chan string
}

var AtomicWorkEventBusInstance *AtomicWorkEventBus

func NewAtomicWorkCenterEventBus() *AtomicWorkEventBus {
	if AtomicWorkEventBusInstance == nil {
		AtomicWorkEventBusInstance = &AtomicWorkEventBus{
			EventBus: make(chan string, 10),
		}
		return AtomicWorkEventBusInstance
	} else {
		return AtomicWorkEventBusInstance
	}
}
