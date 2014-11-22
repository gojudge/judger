package core

import (
	// "encoding/json"
	"fmt"
	"github.com/Unknwon/goconfig"
	// "github.com/duguying/judger/utils"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type Config struct {
	jsonConfigData interface{}
	iniConfigData  *goconfig.ConfigFile
	configType     string
	configFilePath string
}

// load json config file
func (this *Config) loadJsonConfig() error {
	fi, err := os.Open(this.configFilePath)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	configString := string(fd)

	// kick out the comment
	regFilter := regexp.MustCompile(`//[\d\D][^\r]*\r`)
	configString = regFilter.ReplaceAllString(configString, "")
	this.jsonConfigData, err = JsonDecode(configString)
	if err != nil {
		log.Fatalln("Read config file failed. please check `conf/config.json`.")
	}

	return err
}

// load ini config file
func (this *Config) loadIniConfig() error {
	var err error

	this.iniConfigData, err = goconfig.LoadConfigFile(this.configFilePath)
	if nil != err {
		fmt.Println(err)
	}

	return err
}

// get value via key for json
// arg1: key level 1
// arg2: key level 2
func (this *Config) jsonGet(key1 string, key2 string) string {
	// json
	if key1 == "" {
		if value, ok := this.jsonConfigData.(map[string]interface{})[key2].(string); ok {
			return value
		} else {
			fmt.Println(`Cannot Get("", "` + key2 + `")`)
			return ""
		}
	} else {
		if json, ok := this.jsonConfigData.(map[string]interface{})[key1]; ok {
			if value, ok := json.(map[string]interface{})[key2].(string); ok {
				return value
			} else {
				fmt.Println(`Cannot Get("", "` + key2 + `")`)
				return ""
			}
		} else {
			fmt.Println(`Cannot Get("` + key1 + `", "~")`)
			return ""
		}
	}

}

// load config
func (this *Config) LoadConfig(path string) {
	this.configFilePath = path

	regFilter := regexp.MustCompile(`[\d\D][^\r\n]*\.ini$`)
	matched := regFilter.FindAllString(path, -1)

	if 0 != len(matched) {
		this.configType = "ini"
		this.loadIniConfig()
	} else {
		this.configType = "json"
		this.loadJsonConfig()
	}
}

// get value via key
// arg1: key level 1
// arg2: key level 2
func (this *Config) Get(key1 string, key2 string) string {
	if this.configType == "ini" {
		// parse ini
		v, err := this.iniConfigData.GetValue(key1, key2)
		if nil != err {
			fmt.Println(err)
		}

		return v
	} else {
		// parse json
		return this.jsonGet(key1, key2)
	}
}
