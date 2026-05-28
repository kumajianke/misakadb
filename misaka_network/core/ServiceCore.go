package core

import (
	"fmt"
	"misakadb/clilog"
	"misakadb/command"
	"misakadb/config"
	"misakadb/misaka_network"
	onces "misakadb/misaka_network/Onces"
	"misakadb/misaka_network/active"
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

	networkConfig := config.GetGlobalNetworkConfigure()
	maxConn := 1000
	if networkConfig != nil && networkConfig.MaxConn > 0 {
		maxConn = networkConfig.MaxConn
	}
	sem := make(chan struct{}, maxConn)

	for {
		rawConn, err := listener.Accept()
		if err != nil {
			clilog.Warning("accept error:", err)
			continue
		}
		conn := onces.NewSafeConn(rawConn)

		select {
		case sem <- struct{}{}:
			go func(c net.Conn) {
				defer func() { <-sem }()
				serviceCore.handlerConn(c)
			}(conn)
		default:
			clilog.Warning("server full, rejecting connection:", conn.RemoteAddr().String())
			conn.Close()
		}
	}
}

func (serviceCore *ServiceCore) handlerConn(conn net.Conn) {
	defer onces.NewSafeConn(conn).ConnClose()

	connHandler := active.GetServiceConnHandler(conn)

	for {
		client_command, err := connHandler.Recv()
		if err != nil {
			if connHandler.ErrorCounter > 3 {
				onces.NewSafeConn(conn).ConnClose()
				return
			}
			connHandler.ErrorCounter++
			continue
		}

		clilog.Info(fmt.Sprintf(
			"[%s] command: %s",
			conn.RemoteAddr().String(), string(client_command)),
		)

		(command.NewCommandDispatch()).Dispatch(
			connHandler,
			string(client_command),
		)
	}
}
