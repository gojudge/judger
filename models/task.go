package model

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type TaskTab struct {
	Id          int
	TaskId      string
	Language    string
	Type        string
	IoData      string
	Code        string
	Time        time.Time
	BuildLog    string
	BuildResult string
	RunResult   string
	DebugInfo   string
}

func (this *TaskTab) TableName() string {
	return "task"
}

// add task
func (this *TaskTab) AddTask(id int, sid string, language string, ptype string, ioData string, code string) error {
	o := orm.NewOrm()
	var task TaskTab

	task.TaskId = fmt.Sprintf("%s:%d", sid, id)
	task.Language = language
	task.Type = ptype
	task.IoData = ioData
	task.Code = code
	task.Time = time.Now()
	task.BuildLog = ""
	task.BuildResult = ""
	task.RunResult = "TA"
	task.DebugInfo = ""

	_, err = o.Insert(&task)

	return err
}

func (this *TaskTab) GetTaskInfo(id int, sid string) {

}
