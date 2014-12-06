package test

import (
	"github.com/gogather/com"
	"net"
	"testing"
)

func Test_TCP(t *testing.T) {
	addr := "127.0.0.1:1004"

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error("连接服务端失败:", err.Error())
		return
	}

	t.Log("已连接服务器")
	defer conn.Close()

	clientRun(conn)

	conn.Close()
}

func clientRun(conn net.Conn) {
	msg1, _ := com.ReadFileByte("login.json")
	conn.Write(msg1)

	msg2, _ := com.ReadFileByte("task.json")
	conn.Write(msg2)
}
