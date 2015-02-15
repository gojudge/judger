package core

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Filename string
}

func (this *Sqlite) NewSqlite(dataPath string) {
	orm.RegisterDataBase("default", "sqlite3", "data.db")

	this.createTable()
}

func (this *Sqlite) createTable() {
	o := orm.NewOrm()

	_, err := o.Raw("CREATE TABLE [foo] ([id] INTEGER  NOT NULL PRIMARY KEY)").Exec()
	if err != nil {
		fmt.Println(err)
	}

}
