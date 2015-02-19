package model

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type TaskTab struct {
	Id       int
	TaskId   string
	Language string
	Type     string
	IoData   string
	Code     string
	Time     time.Time
}

func (this *TaskTab) TableName() string {
	return "task"
}

func (this *TaskTab) AddTask() {

}
