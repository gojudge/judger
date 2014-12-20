package controller

import (
	"github.com/duguying/judger/core"
	"github.com/duguying/judger/judge"
	"github.com/gogather/com"
	"runtime"
)

// login controller
type LoginController struct {
	core.ControllerInterface
}

func (this *LoginController) Tcp(data map[string]interface{}, cli *core.Client) {
	password := core.C.Get("", "password")

	passwordRecv, ok := data["password"].(string)
	if !ok {
		result, _ := com.JsonEncode(map[string]interface{}{
			"result": false, //bool, login result
			"msg":    "invalid password, password must be string.",
		})
		cli.Write(result)
		cli.Close()
	}

	if password == passwordRecv {
		cli.Login(true)
		result, _ := com.JsonEncode(map[string]interface{}{
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
		result, _ := com.JsonEncode(map[string]interface{}{
			"result": false, //bool, login result
		})
		cli.Write(result)
		cli.Close()
	}

}

// add task controller
type TaskAddController struct {
	core.ControllerInterface
}

func (this *TaskAddController) Tcp(data map[string]interface{}, cli *core.Client) {
	judge.AddTask(data)

	result, _ := com.JsonEncode(map[string]interface{}{
		"result": true, //bool, login result
		"msg":    "response for task",
	})
	cli.Write(result)
}
