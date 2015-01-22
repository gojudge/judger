package controller

import (
	"github.com/duguying/judger/core"
	"github.com/gogather/com"
	"io"
	"net/http"
)

// login controller
type PingController struct {
	core.ControllerInterface
}

func (this *PingController) Tcp(data map[string]interface{}, cli *core.Client) {
	result, _ := com.JsonEncode(map[string]interface{}{
		"result": true,
		"msg":    "pong",
	})
	cli.Write(result)
}

func (this *PingController) Http(data map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	result, _ := com.JsonEncode(map[string]interface{}{
		"result": true,
		"msg":    "pong",
	})

	io.WriteString(w, result)
}
