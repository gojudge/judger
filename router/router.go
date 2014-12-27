package router

import (
	"github.com/duguying/judger/controller"
	"github.com/duguying/judger/core"
)

func init() {
	core.Router("login", &controller.LoginController{})
	core.Router("task_add", &controller.TaskAddController{})
	core.Router("task_info", &controller.GatherController{})
}
