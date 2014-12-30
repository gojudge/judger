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
var Mode string
var configFile string

func Judger() {
	parseArg()

	dataPath := "data.db"

	if Mode == "docker" {
		log.Blueln("[mode]", "docker")

		if !com.FileExist("/data") {
			if err := com.Mkdir("/data"); err != nil {
				log.Warnln("[Warn]", "create dir /data failed")
			} else {
				log.Blueln("[info]", "create dir /data")
			}
		}

		if !com.FileExist("/data/config_docker.ini") {
			com.CopyFile("/data/config_docker.ini", "conf/config_docker.ini")
		}

		if !com.FileExist("/data/executer.json") {
			com.CopyFile("/data/executer.json", "sandbox/c/build/executer.json")
		}

		dataPath = "/data/data.db"
	}

	if configFile == "" {
		configFile = "conf/config.ini"
	}

	if !com.FileExist(configFile) {
		log.Dangerln("[Error]", configFile, "does not exist!")
		os.Exit(-1)
	}

	log.Blueln("[config]")
	log.Blueln(configFile)

	C = &Config{}
	C.NewConfig(configFile)

	GenScript()

	log.Blueln("[data]")
	log.Blueln(dataPath)
	DB = &Sqlite{}
	DB.NewSqlite(dataPath)

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

func parseArg() {
	configFile = ""
	Mode = ""

	arg_num := len(os.Args)

	for i := 0; i < arg_num; i++ {
		s := os.Args[i]

		if s[0] == '-' {
			s = strings.Replace(s, "-", "", -1)
			arr := strings.Split(s, "=")

			if arr[0] == "c" {
				configFile = arr[1]
			} else if arr[0] == "mode" {
				Mode = arr[1]
			}
		}
	}

}
