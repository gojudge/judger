package model

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type SessionTab struct {
	Id         int
	Session    string
	CreateTime time.Time
}

func (this *SessionTab) TableName() string {
	return "session"
}

func (this *SessionTab) CreateSession() {

}

func (this *SessionTab) GetSession(sid string) {

}
