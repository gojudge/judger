package client

import (
	"errors"
	"fmt"
	"github.com/gogather/com/log"
	"io"
	"net"
)

var J *JClient

type JClient struct {
	ip        string
	port      int
	conn      net.Conn
	connected bool
}

func New(ip string, port int) {
	J = &JClient{}
	err := J.Start(ip, port)
	if err != nil {
		J.connected = false
	} else {
		J.connected = true
	}
}

func (this *JClient) Start(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		if conn != nil {
			conn.Close()
		}

		log.Warnln("connect judge server failed in port:", port)
	} else {
		this.conn = conn
		content, _ := this.read()
		log.Blueln(content)
	}

	return err
}

func (this *JClient) Request(msg string) (string, error) {
	if !this.connected {
		return "", errors.New("Not Connected!")
	}
	this.conn.Write([]byte(msg))
	content, err := this.read()
	return content, err
}

func (this *JClient) read() (string, error) {
	var buf [10240]byte

	n, err := this.conn.Read(buf[0:])
	if err != nil {
		if err == io.EOF {
			return "", err
		}
	}

	return string(buf[0:n]), err
}
