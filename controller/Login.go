package controller

import (
	"fmt"
	"github.com/duguying/judger/core"
	"net"
)

type LoginController struct {
	judger.ControllerInterface
}

func (this *LoginController) Tcp(data map[string]interface{}, conn net.Conn) {
	passwordObj := judger.Config("password")
	password, ok := passwordObj.(string)
	if !ok {
		fmt.Println("invalid password in `config.json`, password must be string.")
	}

	passwordRecv, ok := data["password"].(string)
	if !ok {
		fmt.Println("invalid password in `config.json`, password must be string.")
	}

	if password == passwordRecv {
		conn.Write([]byte("this is response from login controller, login success"))
		fmt.Println("response from login controller, login success")
	} else {
		conn.Write([]byte("this is response from login controller, login failed"))
		fmt.Println("response from login controller, login failed")
	}

}
