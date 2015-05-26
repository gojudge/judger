package judge

import (
	"fmt"
	"github.com/gogather/com/log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// run native code(compiled via c/c++) in sandbox
func RunNativeInSandbox(runScript string, runPath string, time int, mem int) error {
	var argTime string
	var argMem string
	var binFilePath string
	var err error

	if time > 0 {
		argTime = fmt.Sprintf("%d", time)
	} else {
		argTime = "10"
	}

	if mem > 0 {
		argMem = fmt.Sprintf("%d", mem)
	} else {
		argMem = "1024"
	}

	currentPath, _ := os.Getwd()
	os.Chdir(runPath)

	if runtime.GOOS == "windows" {
		binFilePath = filepath.Join(runPath, "a.exe")
	} else {
		binFilePath = filepath.Join(runPath, "a.out")
	}

	runScript = filepath.Join(currentPath, runScript)

	log.Pinkln(runScript,
		binFilePath,
		argTime,
		argMem,
	)

	if runtime.GOOS == "windows" {
		err = runnerWin(runScript,
			binFilePath,
			argTime,
			argMem,
		)
	} else {
		err = runnerNix(runScript,
			binFilePath,
			argTime,
			argMem,
		)
	}

	os.Chdir(currentPath)

	return err
}

// call runner in windows
func runnerWin(runScript string, bin string, argTime string, argMem string) error {
	binPath := filepath.Join(bin)
	cmd := exec.Command("cmd", "/K",
		runScript,  // runner script
		argTime+"", // time limit
		argMem,     // memory limit
		binPath,    // executable name
	)

	log.Warnln("[", runScript, binPath, argTime, argMem, "]")

	_, err := cmd.Output()
	if err != nil {
		fmt.Println("失败")
		fmt.Println(err)
		return err
	}

	return nil
}

// call runner in nix
func runnerNix(runScript string, bin string, argTime string, argMem string) error {
	binPath := filepath.Join(bin)
	cmd := exec.Command("sh",
		runScript, // runner script
		argTime,   // time limit
		argMem,    // memory limit
		binPath,   // executable name
	)

	log.Warnln("[", runScript, binPath, argTime, argMem, "]")

	_, err := cmd.Output()
	if err != nil {
		fmt.Println("失败")
		fmt.Println(err)
		return err
	}

	return nil
}
