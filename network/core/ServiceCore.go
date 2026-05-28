package core

import (
	"fmt"
	"misakadb/clilog"
	"misakadb/command"
	"misakadb/network"
	onces "misakadb/network/Onces"
	"misakadb/network/RegisterCenter"
	"misakadb/network/context"
	"net"
	"strconv"
)

type ServiceCore struct {
	ServiceInfo *network.ServiceInfo
}

func NewServiceCore(serviceInfo *network.ServiceInfo) *ServiceCore {
	return &ServiceCore{ServiceInfo: serviceInfo}
}

func (serviceCore *ServiceCore) Close(conn net.Conn) error {
	return conn.Close()
}

func (serviceCore *ServiceCore) Run() error {
	address := serviceCore.ServiceInfo.Address + ":" + strconv.Itoa(serviceCore.ServiceInfo.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		clilog.Error("listen error:", err)
		return err
	}
	defer listener.Close()
	clilog.Success("listening on", address)

	for {
		rawConn, err := listener.Accept()
		if err != nil {
			clilog.Warning("accept error:", err)
			continue
		}
		conn := onces.NewSafeConn(rawConn) // 判断是否为唯一性链接

		registerCenter := RegisterCenter.NewRegisterCenter()
		connContext := context.GetServiceConnContext(conn)

		if registerCenter == nil {
			onces.NewSafeConn(conn).ConnClose()
			clilog.Warning("register center is nil")
			continue
		}

		error := registerCenter.ChanAppendConn(connContext)

		if error != nil {
			clilog.Error("context conn error:", error)
			continue
		}

		go serviceCore.handlerConn(conn, connContext)

	}
}

/**
 * 对接受到的请求进行处理 函数是一个死循环 只有当Exit或者超时的时候会自动断开
 */
func (serviceCore *ServiceCore) handlerConn(conn net.Conn, ConnContext *context.ServiceConnContext) {
	registerCenter := RegisterCenter.NewRegisterCenter()
	defer func() {
		if registerCenter != nil {
			registerCenter.ChanReleaseConn()
		}
		_ = onces.NewSafeConn(conn).ConnClose()
	}()

	for {
		client_command, err := ConnContext.Recv()
		if err != nil {
			if ConnContext.ErrorCounter > 3 {
				return
			}
			ConnContext.ErrorCounter++
			continue
		}

		clilog.Info(fmt.Sprintf(
			"[%s] command: %s",
			conn.RemoteAddr().String(), string(client_command)),
		)

		err = (command.NewCommandDispatch()).Dispatch(
			ConnContext,
			string(client_command),
		)

		if err != nil {
			clilog.Error("dispatch error:", err)
			return
		}
	}
}
