package RegisterCenter

import (
	"errors"
	"fmt"
	"misakadb/clilog"
	onces "misakadb/network/Onces"
	"misakadb/network/context"
	"sync"
)

var lock sync.Mutex

type RegisterCenter struct {
	ConnectQueue chan *context.ServiceConnContext // 链接队列
	Lock         sync.Mutex
}

var RegisterCenterInstance *RegisterCenter = nil

func NewRegisterCenter(connectQueueSize ...int) *RegisterCenter {

	if RegisterCenterInstance != nil {
		return RegisterCenterInstance
	}

	var newConnectQueueSize int
	if len(connectQueueSize) > 0 {
		newConnectQueueSize = connectQueueSize[0]
	} else {
		return nil
	}

	RegisterCenterInstance = &RegisterCenter{
		ConnectQueue: make(chan *context.ServiceConnContext, newConnectQueueSize),
	}

	return RegisterCenterInstance
}

func (connectRegister *RegisterCenter) ChanAppendConn(connContext *context.ServiceConnContext) error {
	select {
	case connectRegister.ConnectQueue <- connContext:
		return nil
	default:
		// 队列已满
		conn := (connContext.Conn)

		err := onces.NewSafeConn(conn).Close()
		if err != nil {
			clilog.Error(fmt.Sprintf("[%s] close conn error", conn.RemoteAddr().String()))
		}
		clilog.Error(fmt.Sprintf("[%s] connect queue is full", conn.RemoteAddr().String()))

		return errors.New("error")
	}
}

func (connectRegister *RegisterCenter) ChanReleaseConn() {
	select {
	case <-connectRegister.ConnectQueue:
	default:
	}
}
