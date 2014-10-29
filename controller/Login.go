package controller

import (
	"fmt"
	"github.com/duguying/judger/core"
	// "net"
)

type LoginController struct {
	judger.ControllerInterface
}

func (this *LoginController) Tcp(data map[string]interface{}, cli *judger.Client) {
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
		cli.Write("response from login controller, login success")
	} else {
		cli.Write("response from login controller, login failed")
		cli.Close()
	}

}
