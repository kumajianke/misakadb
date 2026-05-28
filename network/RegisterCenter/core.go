package RegisterCenter

import (
	"fmt"
	"misakadb/clilog"
	onces "misakadb/network/Onces"
	"misakadb/network/context"
	"sync"
)

var lock sync.Mutex

type ConnectMapperRow struct {
	ConnContext *context.ServiceConnContext
}

type RegisterCenter struct {
	ConnectQueue chan *ConnectMapperRow // 链接队列
	Lock         sync.Mutex
}

var RegisterCenterInstance *RegisterCenter = nil

func NewRegisterCenter(connectQueueSize int) *RegisterCenter {

	if RegisterCenterInstance != nil {
		return RegisterCenterInstance
	}

	RegisterCenterInstance = &RegisterCenter{
		ConnectQueue: make(chan *ConnectMapperRow, connectQueueSize),
	}

	return RegisterCenterInstance
}

func (connectRegister *RegisterCenter) ChanAppendConn(connContext *context.ServiceConnContext) {
	connectMapperRow := &ConnectMapperRow{
		ConnContext: connContext,
	}

	select {
	case connectRegister.ConnectQueue <- connectMapperRow:

	default:
		// 队列已满
		conn := (*connContext.Conn)

		err := onces.NewSafeConn(conn).Close()
		if err != nil {
			clilog.Error(fmt.Sprintf("[%s] close conn error", conn.RemoteAddr().String()))
		}
		clilog.Error(fmt.Sprintf("[%s] connect queue is full", conn.RemoteAddr().String()))
	}
}
