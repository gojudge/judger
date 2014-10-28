package judger

import (
	"github.com/duguying/judger/compiler"
)

func Judger() {
	TcpStart()
	compiler.Compile()
}
