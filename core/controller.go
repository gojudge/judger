package judger

import (
// "net"
)

type ControllerInterface interface {
	Tcp(data map[string]interface{}, cli *Client)
	Http()
}

var RouterMap map[string]ControllerInterface

func Router(actionName string, c ControllerInterface) {
	if nil == RouterMap {
		RouterMap = make(map[string]ControllerInterface)
	}
	RouterMap[actionName] = c
}
