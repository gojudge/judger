package core

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var filename string

type Sqlite struct {
	Filename string
	db       *sql.DB
}

func (this *Sqlite) NewSqlite() {
	var err error
	this.Filename = `data.db`

	this.db, err = sql.Open("sqlite3", filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	if this.checkTableExist("foo") {

	}
}

// check whether the table exist
func (this *Sqlite) checkTableExist(tableName string) bool {
	sql := `SELECT count(*) as num FROM sqlite_master WHERE type='table' AND name='` + tableName + `'`
	fmt.Println(sql)

	rows, err := this.db.Query(sql)
	if err != nil {
		fmt.Println(err)
	}

	var num int
	rows.Scan(&num)
	if num == 0 {
		return false
	} else {
		return true
	}
}

// execute sql
func (this *Sqlite) Exec(sql string) error {
	_, err := this.db.Exec(sql)
	return err
}
