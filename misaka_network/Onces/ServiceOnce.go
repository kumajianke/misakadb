package onces

import (
	"misakadb/clilog"
	"net"
	"sync"
)

var (
	cache_safe_conn       map[net.Conn]*SafeConn = make(map[net.Conn]*SafeConn)
	cache_safe_conn_mutex sync.Mutex
)

type SafeConn struct {
	conn net.Conn
	once sync.Once
}

func (sc *SafeConn) ConnClose() error {
	defer func() {
		cache_safe_conn_mutex.Lock()
		if sc == cache_safe_conn[sc.conn] {
			delete(cache_safe_conn, sc.conn)
		}
		cache_safe_conn_mutex.Unlock()
	}()

	var err error
	sc.once.Do(func() {
		if sc.conn != nil {
			clilog.Info("close conn ", sc.conn.RemoteAddr().String())
			err = sc.conn.Close()
		} else {
			clilog.Warning("get nil conn ")
		}
	})
	return err
}

func NewSafeConn(conn net.Conn) *SafeConn {
	cache_safe_conn_mutex.Lock()
	defer cache_safe_conn_mutex.Unlock()

	if sc, ok := cache_safe_conn[conn]; ok {
		return sc
	} // 缓存命中

	sc := &SafeConn{conn: conn}
	cache_safe_conn[conn] = sc // 回填数据
	return sc
}
