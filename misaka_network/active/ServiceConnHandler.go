package active

import (
	"encoding/binary"
	onces "misakadb/misaka_network/Onces"
	sockshare "misakadb/misaka_network/SockShare"
	"net"
)

type ServiceConnHandler struct {
	Conn         *net.Conn
	ErrorCounter int
}

func getServiceConnHandler(conn net.Conn) *ServiceConnHandler {
	return &ServiceConnHandler{
		Conn:         &conn,
		ErrorCounter: 0,
	}
}

func (handler *ServiceConnHandler) recv() ([]byte, error) {

	conn := *(handler.Conn)

	bytes_lst, err := sockshare.RecvWithHeart(conn)

	if err != nil {
		onces.NewSafeConn(conn).ConnClose()
		return nil, err
	}

	return bytes_lst, nil
}

func (handler *ServiceConnHandler) send(data string) error {
	// 将string转换为[]byte
	dataBytes := []byte(data)

	conn := *(handler.Conn)
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(dataBytes)))
	if _, err := conn.Write(header); err != nil {
		return err
	}
	// 2. 发送数据
	_, err := conn.Write(dataBytes)
	return err
}

func (handler *ServiceConnHandler) commandHandler(command string) {
	_ = command
}
