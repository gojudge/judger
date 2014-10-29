package judger

import (
	"encoding/json"
	"fmt"
	"log"
	// "net"
)

func JsonDecode(data string) (interface{}, error) {
	dataByte := []byte(data)
	var dat interface{}

	err := json.Unmarshal(dataByte, &dat)
	return dat, err
}

func Parse(frame string, cli *Client) {
	fmt.Println(frame)
	json, err := JsonDecode(frame)
	if err != nil {
		log.Print(err)
	} else {
		data := json.(map[string]interface{})

		actonName, ok := data["action"].(string)
		if !ok {
			fmt.Println("invalid request, action name is not exist.")
			cli.conn.Write([]byte(("invalid request, action name is not exist.\n")))
			return
		}

		RouterMap[actonName].Tcp(data, cli)
	}

}
