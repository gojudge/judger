package task

import (
	"encoding/json"
	"fmt"
)

func jsonDecode(data string) interface{} {
	dataByte := []byte(data)
	var dat interface{}

	if err := json.Unmarshal(dataByte, &dat); err != nil {
		panic(err)
	}
	return dat
}

func Run(frame string) {
	fmt.Println(frame)
	json := jsonDecode(frame)
	fmt.Println(json)
}
