package router

import (
	"github.com/duguying/judger/controller"
	"github.com/duguying/judger/core"
)

func init() {
	judger.Router("login", &controller.LoginController{})
}
