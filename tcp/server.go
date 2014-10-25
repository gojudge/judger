package tcp

import (
	"fmt"
	"github.com/duguying/judger/task"
	"io"
	"net"
	"os"
	"regexp"
	// "strings"
	// "time"
)

const (
	BUFF_SIZE = 10
	MARK      = `#`
)

var frame string
var buff = make([]byte, BUFF_SIZE)

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleConnection(tcpConn net.Conn, i int) {
	frame = ""
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
		if n > 0 {
			frame = frame + string(buff[:n])

			reg := regexp.MustCompile(MARK)
			if len(reg.FindAllString(string(buff[:n]), -1)) > 0 {
				// get the json
				frame = reg.ReplaceAllString(frame, "")
				// submit json task
				task.Run(frame)
			}

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
