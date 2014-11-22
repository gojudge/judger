package core

import (
// "fmt"
)

var DB *Sqlite
var C *Config

func Judger() {
	C = &Config{}
	C.NewConfig("conf/config.ini")

	DB = &Sqlite{}
	DB.NewSqlite()

	TcpStart()
}
