package test

import (
	"fmt"
	"github.com/duguying/judger/client"
	"github.com/gogather/com"
	"testing"
)

func Test_TCP(t *testing.T) {
	cli, _ := client.New("127.0.0.1", 1004)

	msg1 := com.ReadFile("login.json")
	resp1, _ := cli.Request(msg1)
	fmt.Println(resp1)

	msg2 := com.ReadFile("task.json")
	resp2, _ := cli.Request(msg2)
	fmt.Println(resp2)

	msg3 := com.ReadFile("info.json")
	resp3, _ := cli.Request(msg3)
	fmt.Println(resp3)

}
