package core

import (
	"net/http"
)

type ControllerInterface interface {
	Tcp(data map[string]interface{}, cli *Client)
	Http(data map[string]interface{}, w http.ResponseWriter, r *http.Request)
}

var RouterMap map[string]ControllerInterface

func Router(actionName string, c ControllerInterface) {
	if nil == RouterMap {
		RouterMap = make(map[string]ControllerInterface)
	}
	RouterMap[actionName] = c
}
