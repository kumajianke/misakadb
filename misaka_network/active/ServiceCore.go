package active

import (
	"fmt"
	"misakadb/clilog"
	"misakadb/misaka_network"
	"net"
	"strconv"
)

type ServiceCore struct {
	ServiceInfo *misaka_network.ServiceInfo
}

func NewServiceCore(serviceInfo *misaka_network.ServiceInfo) *ServiceCore {
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
		conn, err := listener.Accept()
		if err != nil {
			clilog.Warning("accept error:", err)
			continue
		}
		clilog.Info("client connected:", conn.RemoteAddr().String())

		go serviceCore.handlerConn(conn)
	}
}

func (serviceCore *ServiceCore) handlerConn(conn net.Conn) {
	defer serviceCore.Close(conn)

	connHandler := getServiceConnHandler(conn)

	for {
		command, err := connHandler.recv()
		if err != nil {
			if connHandler.ErrorCounter > 3 {
				serviceCore.Close(conn)
				return
			}
			connHandler.ErrorCounter++
			clilog.Error("connHandler recv error:", err)
			continue
		}

		clilog.Info(fmt.Sprintf("[%s] command: %s", conn.RemoteAddr().String(), string(command)))

	}

}
