package controller

import (
	// "fmt"
	"github.com/duguying/judger/core"
	"github.com/gogather/com"
	// "net"
	"github.com/duguying/judger/compiler"
	"html"
	"regexp"
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
	// HTML反转义
	code = html.UnescapeString(code)

	// get the host
	host := cli.Conn.RemoteAddr().String()
	reg := regexp.MustCompile(`:`)
	host = reg.ReplaceAllString(host, "#")

	comp := &compiler.Compile{}
	comp.NewCompile()
	comp.Run(code, language, int(id), host)
}
