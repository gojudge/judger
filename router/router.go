package router

import (
	"github.com/gojudge/judger/controller"
	"github.com/gojudge/judger/core"
)

func init() {
	core.Router("login", &controller.LoginController{})
	core.Router("task_add", &controller.TaskAddController{})
	core.Router("task_info", &controller.GatherController{})
	core.Router("ping", &controller.PingController{})
}
