package compiler

import (
	"fmt"
	"github.com/duguying/judger/core"
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

type Compile struct {
	system      string
	dsm         string // dir split mark
	buildPath   string
	currentPath string

	compiler_c string
	postfix_c  string

	compiler_cpp string
	postfix_cpp  string
}

func (this *Compile) NewCompile() {
	this.system = runtime.GOOS
	this.postfix_c = "c"
	this.postfix_cpp = "cpp"
	this.currentPath, _ = os.Getwd()

	tmpBuildPath := core.C.Get(runtime.GOOS, "buildpath")
	tmpCompilerPath := core.C.Get(runtime.GOOS, "compiler_c")

	regFilter := regexp.MustCompile(`/`)
	if "windows" == runtime.GOOS {
		this.dsm = `\`
		this.buildPath = regFilter.ReplaceAllString(tmpBuildPath, "\\")
		this.compiler_c = regFilter.ReplaceAllString(tmpCompilerPath, "\\")
	} else {
		this.dsm = `/`
		this.buildPath = tmpBuildPath
		this.compiler_c = tmpCompilerPath
	}
}

func (this *Compile) Run(code string, language string, id int, host string) {

	err := this.createDirs(id, host)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		err = this.writeCode(code, id, host, language)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if "windows" == runtime.GOOS {
		this.cl(id, host)
	} else {
		this.gcc(id, host)
	}

}

// 创建编译环境的目录结构
func (this *Compile) createDirs(id int, host string) error {
	var err error
	err = nil
	userBuildPath := this.buildPath + this.dsm + host
	if !core.PathExist(userBuildPath) {
		err = core.Mkdir(userBuildPath)
	}
	itemBuildPath := userBuildPath + this.dsm + fmt.Sprintf("%d", id)
	if !core.PathExist(itemBuildPath) {
		err = core.Mkdir(itemBuildPath)
	}
	return err
}

// 代码写入文件
func (this *Compile) writeCode(code string, id int, host string, language string) error {
	lang := ""
	if language == "C" {
		lang = "c"
	}
	path := this.buildPath + this.dsm + host + this.dsm + fmt.Sprintf("%d%s%d.%s", id, this.dsm, id, lang)
	return core.WriteFile(path, code)
}

// call cl compiler in windows
func (this *Compile) cl(id int, host string) {
	codeFile := this.currentPath + this.dsm + this.buildPath + this.dsm + host + this.dsm + fmt.Sprintf("%d%s%d.c", id, this.dsm, id)
	compiler := this.currentPath + this.dsm + this.compiler_c
	runPath := this.currentPath + this.dsm + this.buildPath + this.dsm + host + this.dsm + fmt.Sprintf("%d", id)

	fmt.Println("codeFile: " + codeFile)
	fmt.Println("compiler: " + compiler)
	fmt.Println("runPath: " + runPath)

	os.Chdir(runPath)

	cmd := exec.Command("cmd", "/K",
		compiler,
		codeFile,
		runPath,
	)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("失败")
		fmt.Println(err)
	}
	fmt.Println(string(output))

	os.Chdir(this.currentPath)
}

// call gcc compiler in other os
func (this *Compile) gcc(id int, host string) {
	codeFile := this.currentPath + this.dsm + this.buildPath + this.dsm + host + this.dsm + fmt.Sprintf("%d%s%d.c", id, this.dsm, id)
	compiler := this.currentPath + this.dsm + this.compiler_c
	runPath := this.currentPath + this.dsm + this.buildPath + this.dsm + host + this.dsm + fmt.Sprintf("%d", id)

	fmt.Println("codeFile: " + codeFile)
	fmt.Println("compiler: " + compiler)
	fmt.Println("runPath: " + runPath)

	os.Chdir(runPath)

	cmd := exec.Command("sh",
		compiler,
		codeFile,
		runPath,
	)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("失败")
		fmt.Println(err)
	}
	fmt.Println(string(output))

	os.Chdir(this.currentPath)
}
