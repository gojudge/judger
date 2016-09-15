package controller

import (
	"github.com/gogather/com"
	"github.com/gojudge/judger/core"
	"github.com/gojudge/judger/judge"
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

func (this *GatherController) Http(data map[string]interface{}, w http.ResponseWriter, r *http.Request) {

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

	io.WriteString(w, result)
}
