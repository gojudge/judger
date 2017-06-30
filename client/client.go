package client

import (
	"errors"
	"fmt"
	"github.com/gogather/com"
	"html"
	"io"
	"log"
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
	sid       string
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

		log.Println("connect judge server failed in port:", port)

		return err
	} else {
		this.conn = conn
		this.connected = true
		content, err := this.read()
		if err != nil {
			return err
		}

		// get seperater mark
		reg := regexp.MustCompile(`/[\d\D]+$`)
		if arr := reg.FindAllString(content, -1); len(arr) > 0 {
			this.mark = com.SubString(arr[0], 1, 1)
		}

		log.Println(content)
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

	sid, ok := response["sid"].(string)
	if ok {
		this.sid = sid
	} else {
		this.sid = ""
	}

	this.login = true

	return err
}

// send request
func (this *JClient) Request(msg map[string]interface{}) (map[string]interface{}, error) {
	msgStr, err := com.JsonEncode(msg)

	if err != nil {
		return nil, err
	}

	if this.conn == nil {
		log.Println("Connection Not Exist")
		return nil, errors.New("Connection Not Exist")
	}

	_, err = this.conn.Write([]byte(msgStr + this.mark))
	if err != nil {
		log.Println("[Write Error]", err)
		this.conn.Close()
		return nil, err
	}

	content, err := this.read()

	if err != nil {
		return nil, err
	}

	// kick sep char
	reg := regexp.MustCompile(this.mark)
	content = reg.ReplaceAllString(content, "")

	if this.debug {
		log.Println("[judger/send:%s]\n%s\n", time.Now(), msgStr)
		log.Println("[judger/recv:%s]\n%s\n", time.Now(), content)
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
func (this *JClient) AddTask(id int64, language string, code string) (map[string]interface{}, error) {
	if !this.login {
		return nil, errors.New("login first")
	}

	req := map[string]interface{}{
		"action":   "task_add",
		"id":       id,
		"sid":      this.sid,
		"time":     time.Now().Nanosecond(),
		"language": language,
		"code":     html.EscapeString(code),
	}

	return this.Request(req)
}

// get task status
func (this *JClient) GetStatus(id int64) (map[string]interface{}, error) {
	if !this.login {
		return nil, errors.New("login first")
	}

	req := map[string]interface{}{
		"action": "task_info",
		"sid":    this.sid,
		"id":     id,
	}

	return this.Request(req)
}

// ping
func (this *JClient) Ping() error {
	if !this.login {
		return errors.New("login first")
	}

	req := map[string]interface{}{
		"action": "ping",
	}

	resp, err := this.Request(req)
	if err != nil {
		return err
	}

	if result, ok := resp["result"].(bool); !ok || !result {
		return errors.New("ping failed.")
	} else {
		return nil
	}
}
