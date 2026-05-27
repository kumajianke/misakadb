package active

import (
	"misakadb/clilog"
	"misakadb/misaka_network"
	"net"
	"strconv"
)

type ServiceCore struct {
	ServiceInfo *misaka_network.ServiceInfo
	connecter   net.Conn
}

func NewServiceCore(serviceInfo *misaka_network.ServiceInfo) *ServiceCore {
	return &ServiceCore{ServiceInfo: serviceInfo}
}

func (serviceCore *ServiceCore) Close() error {
	return serviceCore.connecter.Close()
}

func (serviceCore *ServiceCore) Run() error {
	address := serviceCore.ServiceInfo.Address + ":" + strconv.Itoa(serviceCore.ServiceInfo.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		clilog.Error("listen error:", err)
		return err
	}
	clilog.Success("listening on", address)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			clilog.Warning("accept error:", err)
			continue
		}
		clilog.Info("client connected:", conn.RemoteAddr().String())

		serviceCore.connecter = conn
		go serviceCore.handlerConn(conn)
	}
}

func (serviceCore *ServiceCore) handlerConn(conn net.Conn) {

}
