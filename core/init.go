package core

import (
	"github.com/gogather/com"
	"github.com/gogather/com/log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var DB *Sqlite
var C *Config

func Judger() {
	var configFile string

	if configFile = parseConf(); configFile == "" {
		configFile = "conf/config.ini"
	}

	log.Blueln("[config]")
	log.Blueln(configFile)

	C = &Config{}
	C.NewConfig(configFile)

	GenScript()

	DB = &Sqlite{}
	DB.NewSqlite()

	createBuildDir()

	TcpStart()
}

func createBuildDir() error {
	var err error
	err = nil

	buildPath := filepath.Join(C.Get(runtime.GOOS, "buildpath"))
	if !com.PathExist(buildPath) {
		err = com.Mkdir(buildPath)
	}

	return err
}

func parseConf() string {
	arg_num := len(os.Args)

	for i := 0; i < arg_num; i++ {
		s := os.Args[i]

		if s[0] == '-' {
			s = strings.Replace(s, "-", "", -1)
			arr := strings.Split(s, "=")

			if arr[0] == "c" {
				return arr[1]
			}
		}
	}

	return ""
}
