package onces

import (
	"net"
	"sync"
)

type SafeConn struct {
	net.Conn
	once sync.Once
}

func (sc *SafeConn) Close() error {
	var err error
	sc.once.Do(func() {
		if sc.Conn != nil {
			err = sc.Conn.Close()
		}
	})
	return err
}

func (sc *SafeConn) ConnClose() error {
	return sc.Close()
}

func NewSafeConn(conn net.Conn) *SafeConn {
	if sc, ok := conn.(*SafeConn); ok {
		return sc
	}
	return &SafeConn{Conn: conn}
}
