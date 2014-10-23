package tcp

import (
	"fmt"
	"io"
	"net"
	"os"
	// "time"
)

const BUFF_SIZE = 10

var buff = make([]byte, BUFF_SIZE)

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleConnection(tcpConn net.Conn, i int) {
	if tcpConn == nil {
		return
	}
	for {
		n, err := tcpConn.Read(buff)
		if err == io.EOF {
			fmt.Printf("The RemoteAddr:%s is closed!\n", tcpConn.RemoteAddr().String())
			return
		}
		handleError(err)
		if string(buff[:n]) == "exit" {
			fmt.Printf("The client:%s has exited\n", tcpConn.RemoteAddr().String())
		}
		if n > 0 {
			fmt.Printf("Read:%s\n", string(buff[:n]))
		}
	}
}

func TcpStart() {
	i := 0
	ln, err := net.Listen("tcp", ":1004")
	handleError(err)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			continue
		}
		i += 1
		go handleConnection(conn, i)
	}

}
