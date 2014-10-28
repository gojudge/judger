package controller

import (
	// "fmt"
	"net"
)

func Login(data interface{}, method string, conn net.Conn) {
	conn.Write([]byte("this is response from login controller"))
}
