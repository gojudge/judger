package judger

import (
	"encoding/json"
	"fmt"
	"github.com/duguying/judger/controller"
	"log"
	"net"
)

func jsonDecode(data string) (interface{}, error) {
	dataByte := []byte(data)
	var dat interface{}

	err := json.Unmarshal(dataByte, &dat)
	return dat, err
}

func bind(json interface{}, conn net.Conn) {
	data := json.(map[string]interface{})
	actonName := data["action"].(string)

	if "login" == actonName {
		controller.Login(data, "tcp", conn)
	}
}

func Parse(frame string, conn net.Conn) {
	fmt.Println(frame)
	json, err := jsonDecode(frame)
	if err != nil {
		log.Print(err)
	} else {
		bind(json, conn)
	}

}
