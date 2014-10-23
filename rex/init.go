package judger

import (
	"github.com/duguying/judger/compiler"
	"github.com/duguying/judger/tcp"
)

func Judger() {
	tcp.TcpStart()
	compiler.Compile()
}
