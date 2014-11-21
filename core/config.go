package core

import (
	// "encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var configData interface{}

func readFile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)

	return string(fd)
}

func ConfigInit() {
	var err error
	configString := readFile("conf/config.json")
	// kick out the comment
	regFilter := regexp.MustCompile(`//[\d\D][^\r]*\r`)
	configString = regFilter.ReplaceAllString(configString, "")
	configData, err = JsonDecode(configString)
	if err != nil {
		log.Fatalln("Read config file failed. please check `conf/config.json`.")
	}
}

func Config(key string) interface{} {
	return configData.(map[string]interface{})[key]
}
