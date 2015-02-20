package model

import (
	"github.com/astaxie/beego/orm"
	"github.com/gogather/com"
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

// create session and add session into database
func (this *SessionTab) CreateSession() (string, error) {
	sid := com.CreateGUID()
	o := orm.NewOrm()
	var sess SessionTab

	sess.Session = sid
	sess.CreateTime = time.Now()

	_, err = o.Insert(&sess)

	return sid, err
}

// get session from database
func (this *SessionTab) GetSession(sid string) (Session, error) {
	o := orm.NewOrm()
	var sess SessionTab

	sess.Session = sid
	err = o.Read(&sess, "Session")

	return sess, err
}
