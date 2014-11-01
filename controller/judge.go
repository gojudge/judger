package controller

import (
	"fmt"
	"github.com/duguying/judger/core"
	// "net"
	"github.com/duguying/judger/compiler"
	"runtime"
)

// login controller
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

// add task controller
type TaskAddController struct {
	judger.ControllerInterface
}

func (this *TaskAddController) Tcp(data map[string]interface{}, cli *judger.Client) {

	var ok bool
	var id float64
	var language string
	var code string

	id, ok = data["id"].(float64)
	if !ok {
		cli.Write("invalid id")
		return
	}

	language, ok = data["language"].(string)
	if !ok {
		cli.Write("invalid language name, should be string")
		return
	}

	code, ok = data["code"].(string)
	if !ok {
		cli.Write("invalid code, should be string")
		return
	}

	compiler.Compile(code, language, int(id), "127.0.0.1#5234")
}