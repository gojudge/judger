package core

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"log"
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

	table1 := `CREATE TABLE [session] (
[id] INTEGER  PRIMARY KEY NOT NULL,
[session] VARCHAR(128)  UNIQUE NOT NULL,
[create_time] TIME DEFAULT CURRENT_TIMESTAMP NOT NULL
)`

	table2 := `CREATE TABLE [task] (
[id] INTEGER  PRIMARY KEY AUTOINCREMENT NOT NULL,
[task_id] VARCHAR(128)  UNIQUE NOT NULL,
[language] VARCHAR(64) DEFAULT '''C''' NOT NULL,
[type] VARCHAR(16) DEFAULT '''io''' NOT NULL,
[io_data] TEXT  NULL,
[code] TEXT  NULL,
[time] TIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
[build_log] TEXT  NULL,
[build_result] VARCHAR(128)  NULL,
[run_result] VARCHAR(128)  NULL,
[debug_info] TEXT  NULL
)`

	_, err := o.Raw(table1).Exec()
	if err != nil {
		log.Println(err)
	}

	_, err = o.Raw(table2).Exec()
	if err != nil {
		log.Println(err)
	}

}
