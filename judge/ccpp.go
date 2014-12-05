package judge

import (
	"fmt"
	"github.com/duguying/judger/core"
	"github.com/gogather/com"
	"github.com/gogather/com/log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Compile struct {
	system        string
	buildPath     string
	userBuildPath string
	itemBuildPath string
	codeFilePath  string
	currentPath   string

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

	this.buildPath = filepath.Join(this.currentPath, core.C.Get(runtime.GOOS, "buildpath"))
	this.compiler_c = filepath.Join(this.currentPath, core.C.Get(runtime.GOOS, "compiler_c"))

	log.Blueln("[current path]", this.currentPath)
	log.Blueln("[build path]", this.buildPath)
	log.Blueln("[compiler path]", this.compiler_c)
}

func (this *Compile) Run(code string, language string, id int, sid string) error {

	err := this.createDirs(id, sid)
	if err != nil {
		log.Warnln(err)
		return err
	} else {
		err = this.writeCode(code, id, language)
		if err != nil {
			log.Warnln(err)
			return err
		}
	}

	return this.gcc(id)

}

// 创建编译环境的目录结构
func (this *Compile) createDirs(id int, sid string) error {
	var err error
	err = nil
	this.userBuildPath = filepath.Join(this.buildPath, sid)
	if !com.PathExist(this.userBuildPath) {
		err = com.Mkdir(this.userBuildPath)
	}
	this.itemBuildPath = filepath.Join(this.userBuildPath, fmt.Sprintf("%d", id))
	if !com.PathExist(this.itemBuildPath) {
		err = com.Mkdir(this.itemBuildPath)
	}
	return err
}

// 代码写入文件
func (this *Compile) writeCode(code string, id int, language string) error {
	lang := ""
	if language == "C" {
		lang = "c"
	}
	this.codeFilePath = filepath.Join(this.itemBuildPath, fmt.Sprintf("%d.%s", id, lang))
	return com.WriteFile(this.codeFilePath, code)
}

// call gcc compiler in other os
func (this *Compile) gcc(id int) error {
	os.Chdir(this.itemBuildPath)

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/K",
			this.compiler_c,
			this.codeFilePath,
			this.itemBuildPath,
		)
	} else {
		cmd = exec.Command("sh",
			this.compiler_c,
			this.codeFilePath,
			this.itemBuildPath,
		)
	}

	_, err := cmd.Output()
	if err != nil {
		log.Warnln("失败")
		log.Warnln(err)
	}

	os.Chdir(this.currentPath)

	return err
}
