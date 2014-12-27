package test

import (
	"fmt"
	"github.com/gogather/com"
	"io"
	"net"
	"testing"
)

func Test_TCP(t *testing.T) {
	addr := "127.0.0.1:1004"

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error("连接服务端失败:", err.Error())
		return
	} else {
		fmt.Println("已连接服务器")
		go read(conn)
	}

	defer conn.Close()

	clientRun(conn)

	conn.Close()
}

func clientRun(conn net.Conn) {
	msg1, _ := com.ReadFileByte("login.json")
	conn.Write(msg1)

	read(conn)

	msg2, _ := com.ReadFileByte("task.json")
	conn.Write(msg2)

	read(conn)

	msg3, _ := com.ReadFileByte("info.json")
	conn.Write(msg3)

	read(conn)
}

func read(conn net.Conn) {
	var buf [512]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		if err == io.EOF {
			return
		}
	}
	resp := string(buf[0:n])

	fmt.Printf("[%d]\n", len(resp))
	fmt.Println(resp)
}
