package client

import (
	"fmt"
	"github.com/gogather/com/log"
	"net"
)

type JClient struct {
	ip   string
	port int
}

func (this *JClient) Start(ip string, port int) net.Conn {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		conn.Close()
		log.Warnln("connect judge server failed in port:", port)
		return nil
	} else {
		return conn
	}
}
