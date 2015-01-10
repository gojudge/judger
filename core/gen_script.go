package core

import (
	"fmt"
	"github.com/gogather/com"
	"os"
	"path/filepath"
	"runtime"
)

func GenScript() {
	currentPath, _ := os.Getwd()

	gccWin := `"%s\bin\%s.exe" %%1 -g3 -I"%s\include" -L"%s\lib" -g3 1> BUILD.LOG 2>&1
echo %%ERRORLEVEL%% > BUILDRESULT`

	gccNix := `%s $1 1> BUILD.LOG 2>&1
echo $? > BUILDRESULT`

	var gccScript string
	var gppScript string

	if !com.FileExist("script") {
		com.Mkdir("script")
	}

	if runtime.GOOS == "windows" {
		gccWinPath := C.Get(runtime.GOOS, "gcc_path")
		gccScript = fmt.Sprintf(gccWin, gccWinPath, "gcc", gccWinPath, gccWinPath)
		gppScript = fmt.Sprintf(gccWin, gccWinPath, "g++", gccWinPath, gccWinPath)

		runWin := `"` + filepath.Join(currentPath, C.Get(runtime.GOOS, "executer_path")) + `" %1 %2 %3`

		com.WriteFile(C.Get(runtime.GOOS, "compiler_c"), gccScript)
		com.WriteFile(C.Get(runtime.GOOS, "compiler_cpp"), gppScript)
		com.WriteFile(C.Get(runtime.GOOS, "run_script"), runWin)
	} else {
		gccScript = fmt.Sprintf(gccNix, "gcc")
		gppScript = fmt.Sprintf(gccNix, "g++")
		runNix := filepath.Join(currentPath, C.Get(runtime.GOOS, "executer_path")) + ` $1 $2 $3 -c=` + C.Get(runtime.GOOS, "executer_config")

		com.WriteFile(C.Get(runtime.GOOS, "compiler_c"), gccScript)
		com.WriteFile(C.Get(runtime.GOOS, "compiler_cpp"), gppScript)
		com.WriteFile(C.Get(runtime.GOOS, "run_script"), runNix)

		os.Chmod(C.Get(runtime.GOOS, "compiler_c"), 0755)
		os.Chmod(C.Get(runtime.GOOS, "compiler_cpp"), 0755)
		os.Chmod(C.Get(runtime.GOOS, "run_script"), 0755)
	}

}
