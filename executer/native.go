package executer

import (
	"fmt"
	"os/exec"
)

// run native code(compiled via c/c++) in sandbox
func RunNativeInSandbox(runScript string, bin string, time int, mem int) {
	var argTime string
	var argMem string

	if time > 0 {
		argTime = fmt.Sprintf("-t=%d", time)
	} else {
		argTime = ""
	}

	if mem > 0 {
		argMem = fmt.Sprintf("-m=%d", mem)
	} else {
		argMem = ""
	}

	runnerWin(runScript,
		bin,
		argTime,
		argMem,
	)

}

// call runner in windows
func runnerWin(runScript string, bin string, argTime string, argMem string) error {
	binPath := `.\` + bin + ".exe"
	cmd := exec.Command("cmd", "/K",
		runScript, // runner script
		binPath,   // executable name
		argTime,   // time limit
		argMem,    // memory limit
	)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("失败")
		fmt.Println(err)
		return err
	}
}

// call runner in nix
func runnerNix(runScript string, bin string, argTime string, argMem string) error {
	binPath := `./` + bin
	cmd := exec.Command("sh",
		runScript, // runner script
		bin,       // executable name
		argTime,   // time limit
		argMem,    // memory limit
	)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("失败")
		fmt.Println(err)
		return err
	}
}
