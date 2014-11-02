package compiler

import (
	"fmt"
	"github.com/duguying/judger/core"
	"os/exec"
	"runtime"
)

var buildPath string
var compilerPath string
var DSM string // dir split mark

func Compile(code string, language string, id int, host string) {
	//judger.ConfigInit()
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

  if "windows" == runtime.GOOS {
    DSM = `\`
  }else{
    DSM = `/`
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
		gcc(id, host)
	}
}

// 创建编译环境的目录结构
func createDirs(id int, host string) error {
	var err error
	err = nil
	userBuildPath := buildPath + DSM + host
	if !judger.PathExist(userBuildPath) {
		err = judger.Mkdir(userBuildPath)
	}
	itemBuildPath := userBuildPath + DSM + fmt.Sprintf("%d", id)
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
	path := buildPath + DSM + host + DSM + fmt.Sprintf("%d%s%d.%s", id, DSM, id, lang)
	return judger.WriteFile(path, code)
}

// call cl compiler in windows
func cl(id int, host string) {
	codeFile := buildPath + DSM + host + DSM + fmt.Sprintf("%d%s%d.c", id, DSM, id)

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
func gcc(id int, host string) {
  codeFile := buildPath + DSM + host + DSM + fmt.Sprintf("%d/%d.c", id, id)

	cmd := exec.Command("sh",
		compilerPath,
		codeFile,
    fmt.Sprintf("%s/%s/%d", buildPath, host, id),
	)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("失败")
		fmt.Println(err)
	}
	fmt.Println(string(output))

}

// call g++ compiler
func gpp(id int, host string) {
}
