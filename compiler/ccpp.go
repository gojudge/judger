package compiler

import (
	// "fmt"
	"runtime"
)

func Compile() {
	if "windows" == runtime.GOOS {
		cl()
	} else {
		gcc()
	}
}

// call cl compiler in windows
func cl() {

}

// call gcc compiler in other os
func gcc() {

}
