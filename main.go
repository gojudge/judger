package main

import (
	"github.com/duguying/judger/compiler"
	// "github.com/duguying/judger/core"
	_ "github.com/duguying/judger/router"
)

func main() {
	code := `
#include <stdio.h>

int main(void){
  printf("hello world.\n");
  return 0;
}
`
	compiler.Compile(code, "C", 2815, "127.0.0.1#5234")
	// judger.Judger()
}
