package eventbus

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
