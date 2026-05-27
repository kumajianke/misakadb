package active

import (
	"encoding/binary"
	"io"
	"misakadb/clilog"
	"net"
)

type ServiceConnHandler struct {
	Conn *net.Conn
}

func (handler *ServiceConnHandler) recv() []byte {
	conn := *(handler.Conn)

	len_recv := make([]byte, 4)
	_, err_len := io.ReadFull(conn, len_recv)
	if err_len != nil {
		clilog.Error("bad recv from conn" + conn.RemoteAddr().String())
		return nil
	}
	len_number := binary.BigEndian.Uint32(len_recv)

	bytes_lst := make([]byte, len_number)
	_, err := io.ReadFull(conn, bytes_lst)

	if err != nil {
		clilog.Error("bad recv from conn" + conn.RemoteAddr().String())
		return nil

	}
	data := bytes_lst[:len_number]
	return data
}
