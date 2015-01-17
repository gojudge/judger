package controller

import (
	"github.com/duguying/judger/core"
	"github.com/duguying/judger/judge"
	"github.com/gogather/com"
	"io"
	"net/http"
	"runtime"
)

// gather information controller
type GatherController struct {
	core.ControllerInterface
}

func (this *GatherController) Tcp(data map[string]interface{}, cli *core.Client) {
	info := &judge.Info{}
	sid := data["sid"].(string)
	id := data["id"].(float64)
	information := info.Gather(sid, int(id), core.C.Get(runtime.GOOS, "buildpath"))

	result, _ := com.JsonEncode(map[string]interface{}{
		"info": information,

		"time": 123456789,
		"sid":  sid,
		"id":   id,
	})
	cli.Write(result)
}

func (this *GatherController) Http(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello, world from controller\n")
}
