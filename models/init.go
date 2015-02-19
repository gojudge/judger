package models

import (
	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(SessionTab))
	orm.RegisterModel(new(TaskTab))
}
