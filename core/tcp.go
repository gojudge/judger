package judger

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
)

const (
	BUFF_SIZE = 10
	MARK      = `#`
	MAX_LCI   = 100
)

type Client struct {
	active bool
	conn   net.Conn
}

var buff = make([]byte, BUFF_SIZE)
var cliTab = make(map[int]*Client)

/// close client connect from server
func (this *Client) Close() {
	this.conn.Close()
	this.active = false
}

// send message to client and print in server console
func (this *Client) Write(str string) {
	this.conn.Write([]byte(str))
	fmt.Println(str)
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func handleConnection(tcpConn net.Conn, cid int) {
	frame := ""

	if tcpConn == nil {
		return
	}

	cli := &Client{true, tcpConn}
	cliTab[cid] = cli

	fmt.Println("Connected! Remote address is " + tcpConn.LocalAddr().String())
	tcpConn.Write([]byte("Connected! Remote address is " + tcpConn.LocalAddr().String()))
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
				// kick out the comment
				regFilter := regexp.MustCompile(`//[\d\D][^\r]*\r`)
				frame = regFilter.ReplaceAllString(frame, "")
				// get the json
				frame = reg.ReplaceAllString(frame, "")
				// submit json task
				Parse(frame, cli)
				frame = ""
				// if connection is inactive[closed by server, jump out of cycle
				if !cli.active {
					return
				}
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
		if i > MAX_LCI {
			fmt.Println("reached max client limit, server stoped.")
			return
		}
		go handleConnection(conn, i)
	}

}
