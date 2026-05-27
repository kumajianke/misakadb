package active

import (
	"encoding/binary"
	"errors"
	"io"
	"misakadb/clilog"
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

	len_recv := make([]byte, 4)
	_, err_len := io.ReadFull(conn, len_recv)
	if err_len != nil {
		clilog.Error("bad recv from conn " + conn.RemoteAddr().String())
		return nil, errors.New("bad recv from conn" + err_len.Error())

	}
	len_number := binary.BigEndian.Uint32(len_recv)

	bytes_lst := make([]byte, len_number)
	_, err := io.ReadFull(conn, bytes_lst)

	if err != nil {
		clilog.Error("bad recv from conn " + conn.RemoteAddr().String())
		return nil, errors.New("bad recv from conn" + err.Error())

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
