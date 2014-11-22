package core

import (
// "fmt"
)

var DB *Sqlite
var C *Config

func Judger() {
	C = &Config{}
	C.LoadConfig("conf/config.ini")

	// pwd := C.Get("", "password")
	// fmt.Println("password: " + pwd)

	DB = &Sqlite{}
	DB.NewSqlite()

	TcpStart()
}
