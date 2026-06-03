package RegisterCenter

import (
	"errors"
	"fmt"
	"misakadb/clilog"
	onces "misakadb/network/Onces"
	"misakadb/network/context"
	"os"
	"sync"
)

var lock sync.Mutex

type RegisterCenter struct {
	ConnectQueue chan *context.ServiceConnContext // 链接队列
	MasterKey    string                           // 密钥
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
		// 无参数调用时，创建一个用于工具/加密场景的简化实例（不初始化连接队列）
		newConnectQueueSize = 0
	}

	key, errors := os.ReadFile("./profiles/master.mikey")

	if errors != nil {
		clilog.Error("[严重错误] 无法获取到密钥!")
		panic("service is not runnable")
	}

	var connQueue chan *context.ServiceConnContext
	if newConnectQueueSize > 0 {
		connQueue = make(chan *context.ServiceConnContext, newConnectQueueSize)
	}

	RegisterCenterInstance = &RegisterCenter{
		ConnectQueue: connQueue,
		MasterKey:    string(key),
	}

	return RegisterCenterInstance
}

func (connectRegister *RegisterCenter) ChanAppendConn(connContext *context.ServiceConnContext) error {
	// 检查队列是否已初始化
	if connectRegister.ConnectQueue == nil {
		conn := connContext.Conn
		err := onces.NewSafeConn(conn).Close()
		if err != nil {
			clilog.Error(fmt.Sprintf("[%s] close conn error", conn.RemoteAddr().String()))
		}
		return errors.New("connect queue not initialized")
	}

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
	// 检查队列是否已初始化
	if connectRegister.ConnectQueue == nil {
		return
	}

	select {
	case <-connectRegister.ConnectQueue:
	default:
	}
}
