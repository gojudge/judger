package judger

import (
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"os"
)

func JsonDecode(data string) (interface{}, error) {
	dataByte := []byte(data)
	var dat interface{}

	err := json.Unmarshal(dataByte, &dat)
	return dat, err
}

func JsonEncode(data interface{}) string {
	a, _ := json.Marshal(data)
	return string(a)
}

// 检查文件或目录是否存在
// 如果由 pathname 指定的文件或目录存在则返回 true，否则返回 false
func PathExist(pathname string) bool {
	_, err := os.Stat(pathname)
	return err == nil || os.IsExist(err)
}

// 创建文件夹
func Mkdir(path string) error {
	return os.Mkdir(path, os.ModePerm)
}

// 字符串写入文件
func WriteFile(fullpath string, str string) error {
	data := []byte(str)
	return ioutil.WriteFile(fullpath, data, 0644)
}
