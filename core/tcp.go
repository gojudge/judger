package judger

import (
	"fmt"
	"io"
	"log"
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
	Conn   net.Conn
	cid    int
	login  bool
}

var buff = make([]byte, BUFF_SIZE)
var cliTab = make(map[int]*Client)

/// close client connect from server
func (this *Client) Close() {
	this.Conn.Close()
	this.active = false
	cliTab[this.cid] = nil
}

// send message to client and print in server console
func (this *Client) Write(str string) {
	this.Conn.Write([]byte(str))
	fmt.Println(str)
}

// set mark for login
// value: true for login, false for not login
func (this *Client) Login(value bool) {
	this.login = value
}

func Parse(frame string, cli *Client) {
	fmt.Println(frame)
	json, err := JsonDecode(frame)
	if err != nil {
		log.Print(err)
	} else {
		data := json.(map[string]interface{})

		actonName, ok := data["action"].(string)
		if !ok {
			cli.Write("invalid request, action name is not exist.")
			return
		}

		// 如果不是登录请求，并且用户处于未登录状态，禁止通行
		if actonName != "login" {
			if !cli.login {
				cli.Write("you have not login.")
				return
			}
		}

		RouterMap[actonName].Tcp(data, cli)
	}

}

func handleError(err error, tcpConn net.Conn) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Client error: %s", err.Error())
		if tcpConn != nil {
			tcpConn.Close()
		}
	}
}

func handleConnection(tcpConn net.Conn, cid int) {
	frame := ""

	if tcpConn == nil {
		return
	}

	cli := &Client{true, tcpConn, cid, false}
	cliTab[cid] = cli

	fmt.Println("Connected! Remote address is " + tcpConn.LocalAddr().String())
	tcpConn.Write([]byte("Connected! Remote address is " + tcpConn.LocalAddr().String()))
	for {
		n, err := tcpConn.Read(buff)
		if err == io.EOF {
			fmt.Printf("The RemoteAddr:%s is closed!\n", tcpConn.RemoteAddr().String())
			return
		}
		handleError(err, tcpConn)
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
	handleError(err, nil)

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
