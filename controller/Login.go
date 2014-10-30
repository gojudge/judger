package controller

import (
	"fmt"
	"github.com/duguying/judger/core"
	// "net"
	"runtime"
)

type LoginController struct {
	judger.ControllerInterface
}

func (this *LoginController) Tcp(data map[string]interface{}, cli *judger.Client) {
	passwordObj := judger.Config("password")
	password, ok := passwordObj.(string)

	if !ok {
		fmt.Println("invalid password in `config.json`, password must be string.")
		result := judger.JsonEncode(map[string]interface{}{
			"result": false, //bool, login result
			"msg":    "internal error!",
		})
		cli.Write(result)
		cli.Close()
	}

	passwordRecv, ok := data["password"].(string)
	if !ok {
		result := judger.JsonEncode(map[string]interface{}{
			"result": false, //bool, login result
			"msg":    "invalid password, password must be string.",
		})
		cli.Write(result)
		cli.Close()
	}

	if password == passwordRecv {
		cli.Login(true)
		result := judger.JsonEncode(map[string]interface{}{
			"result": true, //bool, login result
			"os":     runtime.GOOS + " " + runtime.GOARCH,
			"language": map[string]interface{}{ //language:compiler
				"C":    "gcc",
				"C++":  "g++",
				"Java": "javac version 1.7",
			},
			"time": 123456789, //server time stamp
		})
		cli.Write(result)
	} else {
		result := judger.JsonEncode(map[string]interface{}{
			"result": false, //bool, login result
		})
		cli.Write(result)
		cli.Close()
	}

}
