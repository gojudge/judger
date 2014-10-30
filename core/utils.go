package judger

import (
	"encoding/json"
	// "fmt"
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
