package context

import (
	"encoding/binary"
	onces "misakadb/network/Onces"
	sockshare "misakadb/network/SockShare"
	"net"
)

/**
 * 服务器中单个连接上下文，每个连接都有一个上下文, 记录连接的状态及其行动操作
 */
type ServiceConnContext struct {
	Conn         net.Conn
	ErrorCounter int
	Detail       map[string]any
}

func GetServiceConnContext(conn net.Conn) *ServiceConnContext {
	return &ServiceConnContext{
		Conn:         conn,
		ErrorCounter: 0,
		Detail:       make(map[string]any),
	}
}

func (context *ServiceConnContext) Recv() ([]byte, error) {

	conn := context.Conn

	bytes_lst, err := sockshare.RecvWithHeart(conn)

	if err != nil {
		onces.NewSafeConn(conn).ConnClose()
		return nil, err
	}

	return bytes_lst, nil
}

func (context *ServiceConnContext) Send(data string) error {
	// 将string转换为[]byte
	dataBytes := []byte(data)

	conn := context.Conn
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(dataBytes)))
	if _, err := conn.Write(header); err != nil {
		return err
	}
	// 2. 发送数据
	_, err := conn.Write(dataBytes)
	return err
}
