package compiler

import (
	"fmt"
	"github.com/duguying/judger/core"
	"os/exec"
	"runtime"
)

var buildPath string
var compilerPath string

func Compile(code string, language string, id int, host string) {
	judger.ConfigInit()
	var ok bool
	buildPathObj := judger.Config("buildpath")
	buildPath, ok = buildPathObj.(string)
	if !ok {
		fmt.Println("`buildpath` is error in config.json")
	}

	compilerPathObj := judger.Config("compilerpath")
	compilerPath, ok = compilerPathObj.(string)
	if !ok {
		fmt.Println("`compilerpath` is error in config.json")
	}

	err := createDirs(id, host)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		err = writeCode(code, id, host, language)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if "windows" == runtime.GOOS {
		cl(id, host)
	} else {
		gcc()
	}
}

// 创建编译环境的目录结构
func createDirs(id int, host string) error {
	var err error
	err = nil
	userBuildPath := buildPath + `\` + host
	if !judger.PathExist(userBuildPath) {
		err = judger.Mkdir(userBuildPath)
	}
	itemBuildPath := userBuildPath + `\` + fmt.Sprintf("%d", id)
	if !judger.PathExist(itemBuildPath) {
		err = judger.Mkdir(itemBuildPath)
	}
	return err
}

// 代码写入文件
func writeCode(code string, id int, host string, language string) error {
	lang := ""
	if language == "C" {
		lang = "c"
	}
	path := buildPath + `\` + host + `\` + fmt.Sprintf("%d\\%d.%s", id, id, lang)
	return judger.WriteFile(path, code)
}

// call cl compiler in windows
func cl(id int, host string) {
	codeFile := buildPath + `\` + host + `\` + fmt.Sprintf("%d\\%d.c", id, id)

	cmd := exec.Command("cmd", "/K",
		compilerPath,
		codeFile,
		fmt.Sprintf("%s\\%s\\%d", buildPath, host, id),
	)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("失败")
		fmt.Println(err)
	}
	fmt.Println(string(output))
}

// call gcc compiler in other os
func gcc() {

}
