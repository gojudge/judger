package test

import (
	"fmt"
	"github.com/duguying/judger/client"
	"testing"
)

func Test_TCP(t *testing.T) {
	cli, err := client.New("duguying.net", 1004, "123456789")

	if err != nil {
		fmt.Println(err)
	}

	response, err := cli.AddTask(12, "randomstring", "C", "int main(){return 0;}")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}

	response, err = cli.GetStatus(12, "randomstring")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}

}
