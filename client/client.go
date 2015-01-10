package client

import (
	"fmt"
	"github.com/gogather/com/log"
	"net"
)

type JClient struct {
	ip   string
	port int
	conn net.Conn
}

func (this *JClient) Start(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		conn.Close()
		log.Warnln("connect judge server failed in port:", port)
	} else {
		this.conn = conn
	}

	return err
}

func (this *JClient) Request(msg string) string {
	var buf [10240]byte

	this.conn.Write(msg)

	n, err := this.conn.Read(buf[0:])
	if err != nil {
		if err == io.EOF {
			return
		}
	}

	return string(buf[0:n])
}
