package core

var DB Sqlite

func Judger() {
	ConfigInit()
	TcpStart()

	DB := &Sqlite{}
	DB.NewSqlite()
}
