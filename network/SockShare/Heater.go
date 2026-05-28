package sockshare

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"misakadb/clilog"
	"misakadb/config"
	"net"
	"os"
	"time"
)

/**
* 通过心跳检测的方式接受数据
* @param conn 连接对象
* @param recvLength 接收长度
 */
func RecvWithHeart(conn net.Conn) ([]byte, error) {
	networkConfig := config.GetGlobalNetworkConfigure()

	retryCount := networkConfig.RetryCount // 心跳重试次数
	retryDelay := networkConfig.RetryDelay // 心跳重试延迟

	errorRecvCounter := 0

	for {
		// 设置心跳超时
		conn.SetDeadline(
			time.Now().Add(time.Duration(retryDelay) * time.Second),
		)

		// 接受长度
		recv_len := make([]byte, 4)
		_, err_len := io.ReadFull(conn, recv_len)

		if err_len != nil {
			// 如果错误是超时
			if errors.Is(err_len, os.ErrDeadlineExceeded) {
				errorRecvCounter++
				if errorRecvCounter > retryCount {
					return nil, errors.New("can not recv any data")
				}
				continue
			}
			// 其他错误
			return nil, errors.New("bad recv step 1 from conn " + err_len.Error())
		}
		len_number := binary.BigEndian.Uint32(recv_len)
		errorRecvCounter = 0 // 成功就将重试次数清空

		// 最大载荷检查避免OOM
		const MaxPayloadSize = 16 * 1024 * 1024 // 16MB
		if len_number > MaxPayloadSize {
			return nil, errors.New("payload too large")
		}

		// 获取输入的类型
		recv_type := make([]byte, 1)
		_, err_type := io.ReadFull(conn, recv_type)
		if err_type != nil {
			return nil, errors.New("bad recv step 2 from conn " + err_type.Error())
		}

		if recv_type[0] == 0x01 {
			// 这个只是我们的心跳包 我们继续等待数据包
			clilog.Info(fmt.Sprintf("[%s] get heartbeat of conn", conn.RemoteAddr().String()))
			continue
		}

		// 获取数据
		bytes_lst := make([]byte, len_number)
		_, err := io.ReadFull(conn, bytes_lst)

		if err != nil {
			return nil, errors.New("bad recv from conn " + err.Error())
		}

		return bytes_lst, nil

	}
}
