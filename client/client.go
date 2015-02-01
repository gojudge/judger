package client

import (
	"errors"
	"fmt"
	"github.com/gogather/com"
	"github.com/gogather/com/log"
	"io"
	"net"
	"regexp"
	"time"
)

var J *JClient

type JClient struct {
	ip        string
	port      int
	conn      net.Conn
	connected bool
	mark      string
	debug     bool
}

func New(ip string, port int) (*JClient, error) {
	J = &JClient{}
	err := J.Start(ip, port)
	if err != nil {
		J.connected = false
	} else {
		J.connected = true
	}
	return J, err
}

func (this *JClient) SetDebug(flag bool) {
	this.debug = flag
}

func (this *JClient) Start(ip string, port int) error {
	// default # to get the real sep
	this.mark = "#"
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
		// get seperater mark
		reg := regexp.MustCompile(`/[\d\D]+$`)
		if arr := reg.FindAllString(content, -1); len(arr) > 0 {
			this.mark = com.SubString(arr[0], 1, 1)
		}

		log.Blueln(content)
	}

	return err
}

func (this *JClient) Request(msg string) (string, error) {
	if !this.connected {
		return "", errors.New("Not Connected!")
	}
	this.conn.Write([]byte(msg + this.mark))
	content, err := this.read()
	// kick sep char
	reg := regexp.MustCompile(this.mark)
	content = reg.ReplaceAllString(content, "")

	if this.debug {
		log.Bluef("[judger/send:%s]\n%s\n", time.Now(), msg)

		log.Warnf("[judger/recv:%s]\n%s\n", time.Now(), content)
	}

	return content, err
}

func (this *JClient) read() (string, error) {
	var buff [10]byte
	frame := ""

	for {
		n, err := this.conn.Read(buff[0:])
		if err != nil {
			if err == io.EOF {
				return "", err
			}
		}

		if n > 0 {
			frame = frame + string(buff[:n])

			reg := regexp.MustCompile(this.mark)
			if len(reg.FindAllString(string(buff[:n]), -1)) > 0 {
				break
			}

		}
	}

	return frame, nil
}
