package controller

import (
	"github.com/duguying/judger/core"
	"github.com/duguying/judger/judge"
	"github.com/duguying/judger/models"
	"github.com/gogather/com"
	"io"
	"net/http"
	"runtime"
)

// login controller
type LoginController struct {
	core.ControllerInterface
}

func (this *LoginController) Tcp(data map[string]interface{}, cli *core.Client) {
	sess := &models.SessionTab{}
	password := core.C.Get("", "password")

	passwordRecv, ok := data["password"].(string)
	if !ok {
		result, _ := com.JsonEncode(map[string]interface{}{
			"result": false, //bool, login result
			"msg":    "invalid password, password must be string.",
		})
		cli.Write(result)
		return
	}
	if password != passwordRecv {
		result, _ := com.JsonEncode(map[string]interface{}{
			"result": false, //bool, login result
			"msg":    "wrong password.",
		})
		cli.Write(result)
		cli.Close()
		return
	}

	sid, ok := data["sid"]
	if ok {
		sidString, ok := sid.(string)
		if !ok {
			result, _ := com.JsonEncode(map[string]interface{}{
				"result": false, //bool, login result
				"msg":    "invalid sid, sid must be a string.",
			})
			cli.Write(result)
			return
		}

		s, err := sess.GetSession(sidString)
		if err == nil && s.Session == sidString {
			// login success
			cli.Login(true)
			result, _ := com.JsonEncode(map[string]interface{}{
				"result": true, //bool, login result
				"sid":    sid,
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
				"msg":    "invalid sid, please login without sid.",
			})
			cli.Write(result)
		}
	} else {
		sid, err := sess.CreateSession()
		if err != nil {
			result, _ := com.JsonEncode(map[string]interface{}{
				"result": false, //bool, login result
				"msg":    "create session failed. please contact the admin.",
			})
			cli.Write(result)
		} else {
			cli.Login(true)
			result, _ := com.JsonEncode(map[string]interface{}{
				"result": true, //bool, login result
				"sid":    sid,
				"os":     runtime.GOOS + " " + runtime.GOARCH,
				"language": map[string]interface{}{ //language:compiler
					"C":    "gcc",
					"C++":  "g++",
					"Java": "javac version 1.7",
				},
				"time": 123456789, //server time stamp
			})
			cli.Write(result)
		}
	}

}

func (this *LoginController) Http(data map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	var sid string
	sess := &models.SessionTab{}

	sidObj, ok := data["sid"]
	if !ok {
		sid, _ = sess.CreateSession()
	} else {
		sid, ok = sidObj.(string)
		if !ok {
			result, _ := com.JsonEncode(map[string]interface{}{
				"result": false, //bool, login result
				"msg":    "invalid sid, sid must be a string.",
			})
			io.WriteString(w, result)
			return
		} else {
			s, err := sess.GetSession(sid)
			if err != nil {
				result, _ := com.JsonEncode(map[string]interface{}{
					"result": false, //bool, login result
					"msg":    "invalid sid, login failed.",
				})
				io.WriteString(w, result)
				return
			}

			if s.Session == sid {

			} else {
				result, _ := com.JsonEncode(map[string]interface{}{
					"result": false, //bool, login result
					"msg":    "invalid sid, please login without sid.",
				})
				io.WriteString(w, result)
				return
			}
		}
	}

	// login success
	result, _ := com.JsonEncode(map[string]interface{}{
		"result": true, //bool, login result
		"sid":    sid,
		"os":     runtime.GOOS + " " + runtime.GOARCH,
		"language": map[string]interface{}{ //language:compiler
			"C":    "gcc",
			"C++":  "g++",
			"Java": "javac version 1.7",
		},
		"time": 123456789, //server time stamp
	})
	io.WriteString(w, result)

}

// add task controller
type TaskAddController struct {
	core.ControllerInterface
}

func (this *TaskAddController) Tcp(data map[string]interface{}, cli *core.Client) {
	judge.AddTask(data)

	result, _ := com.JsonEncode(map[string]interface{}{
		"result": true, //bool, login result
		"msg":    "task added",
	})
	cli.Write(result)
}

func (this *TaskAddController) Http(data map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	judge.AddTask(data)

	result, _ := com.JsonEncode(map[string]interface{}{
		"result": true, //bool, login result
		"msg":    "task added",
	})
	io.WriteString(w, result)
}
