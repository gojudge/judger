package client

import (
	"errors"
	"fmt"
	"github.com/gogather/com"
	"github.com/gogather/com/log"
	"html"
	"io"
	"net"
	"regexp"
	"time"
)

type JClient struct {
	ip        string
	port      int
	conn      net.Conn
	connected bool
	mark      string
	debug     bool
	login     bool
}

func New(ip string, port int, password string) (*JClient, error) {
	J := &JClient{}
	err := J.Start(ip, port, password)
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

// start the session
func (this *JClient) Start(ip string, port int, password string) error {
	// not login, first time
	this.login = false

	// default # to get the real sep
	this.mark = "#"
	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		if conn != nil {
			conn.Close()
		}

		log.Warnln("connect judge server failed in port:", port)

		return err
	} else {
		this.conn = conn
		content, err := this.read()
		if err != nil {
			return err
		}

		// get seperater mark
		reg := regexp.MustCompile(`/[\d\D]+$`)
		if arr := reg.FindAllString(content, -1); len(arr) > 0 {
			this.mark = com.SubString(arr[0], 1, 1)
		}

		log.Blueln(content)
	}

	// login
	loginRequest := map[string]interface{}{
		"action":   "login",
		"password": password,
	}

	response, err := this.Request(loginRequest)
	if err != nil {
		return err
	}

	result, ok := response["result"].(bool)
	if !result || !ok {
		return errors.New("login failed.")
	}

	this.login = true

	return err
}

// send request
func (this *JClient) Request(msg map[string]interface{}) (map[string]interface{}, error) {
	if !this.connected {
		return nil, errors.New("Not Connected!")
	}

	msgStr, err := com.JsonEncode(msg)

	if err != nil {
		return nil, err
	}

	this.conn.Write([]byte(msgStr + this.mark))
	content, err := this.read()

	if err != nil {
		return nil, err
	}

	// kick sep char
	reg := regexp.MustCompile(this.mark)
	content = reg.ReplaceAllString(content, "")

	if this.debug {
		log.Bluef("[judger/send:%s]\n%s\n", time.Now(), msgStr)
		log.Warnf("[judger/recv:%s]\n%s\n", time.Now(), content)
	}

	resp, err := com.JsonDecode(content)

	return resp.(map[string]interface{}), err
}

// read message from socket
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

// add task
func (this *JClient) AddTask(id int64, sid string, language string, code string) (map[string]interface{}, error) {
	if !this.login {
		return nil, errors.New("login first")
	}

	req := map[string]interface{}{
		"action":   "task_add",
		"id":       id,
		"sid":      sid,
		"time":     time.Now().Nanosecond(),
		"language": language,
		"code":     html.EscapeString(code),
	}

	return this.Request(req)
}

// get task status
func (this *JClient) GetStatus(id int64, sid string) (map[string]interface{}, error) {
	if !this.login {
		return nil, errors.New("login first")
	}

	req := map[string]interface{}{
		"action": "task_info",
		"sid":    sid,
		"id":     id,
	}

	return this.Request(req)
}
