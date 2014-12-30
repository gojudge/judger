package core

import (
	"fmt"
	"github.com/gogather/com"
	"os"
	"runtime"
)

func GenScript() {
	currentPath, _ := os.Getwd()

	gccWin := `@set PATH=%%PATH%%;%s
%s %%1 1> BUILD.LOG 2>&1
echo %%ERRORLEVEL%% > BUILDRESULT`

	gccNix := `%s $1 > BUILD.LOG
echo $? > BUILDRESULT`

	var gccScript string
	var gppScript string

	if !com.FileExist("script") {
		com.Mkdir("script")
	}

	if runtime.GOOS == "windows" {
		gccWinPath := C.Get(runtime.GOOS, "gcc_path")
		gccScript = fmt.Sprintf(gccWin, gccWinPath, "gcc")
		gppScript = fmt.Sprintf(gccWin, gccWinPath, "g++")
		runWin := `"` + currentPath + `\sandbox\c\build\executer.exe" %1 %2 %3`

		com.WriteFile(C.Get(runtime.GOOS, "compiler_c"), gccScript)
		com.WriteFile(C.Get(runtime.GOOS, "compiler_cpp"), gppScript)
		com.WriteFile(C.Get(runtime.GOOS, "run_script"), runWin)
	} else {
		gccScript = fmt.Sprintf(gccNix, "gcc")
		gppScript = fmt.Sprintf(gccNix, "g++")
		runNix := currentPath + `/sandbox/c/build/executer %1 %2 %3 -c=` + C.Get(runtime.GOOS, "executer_config")

		com.WriteFile(C.Get(runtime.GOOS, "compiler_c"), gccScript)
		com.WriteFile(C.Get(runtime.GOOS, "compiler_cpp"), gppScript)
		com.WriteFile(C.Get(runtime.GOOS, "run_script"), runNix)

		os.Chmod(C.Get(runtime.GOOS, "compiler_c"), 0755)
		os.Chmod(C.Get(runtime.GOOS, "compiler_cpp"), 0755)
		os.Chmod(C.Get(runtime.GOOS, "run_script"), 0755)
	}

}
