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

	connectQueue := RegisterCenter.RegisterCenterInstance.ConnectQueue
	for {
		rawConn, err := listener.Accept()
		if err != nil {
			clilog.Warning("accept error:", err)
			continue
		}
		conn := onces.NewSafeConn(rawConn)

		connContext := context.GetServiceConnContext(conn)
		connectMapperRow := &RegisterCenter.ConnectMapperRow{
			ConnContext: connContext,
		}

		select {
		case connectQueue <- connectMapperRow:
			serviceCore.contextConn(conn, connContext)
			clilog.Success("connect success:", conn.RemoteAddr().String())
		default:
			clilog.Warning("server full, rejecting connection:", conn.RemoteAddr().String())
			conn.Close()
			conn = nil
		}
	}
}

/**
 * 处理连接用的上下文
 */
func (serviceCore *ServiceCore) contextConn(conn net.Conn, ConnContext *context.ServiceConnContext) {
	defer onces.NewSafeConn(conn).ConnClose()

	for {
		client_command, err := ConnContext.Recv()
		if err != nil {
			if ConnContext.ErrorCounter > 3 {
				onces.NewSafeConn(conn).ConnClose()
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
