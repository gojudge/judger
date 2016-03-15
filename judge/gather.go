package judge

import (
	"fmt"
	"github.com/gogather/com"
	"github.com/gogather/com/log"
	"path/filepath"
	"strconv"
)

type Info struct {
	sid         string
	id          int
	buildLog    string
	buildResult int
	runResult   string
	buildPath   string
}

func (this *Info) Gather(sid string, id int, buildPath string) map[string]interface{} {
	this.sid = sid
	this.id = id
	this.buildPath = buildPath

	this.buildLog = this.getLog("BUILD.LOG")

	if this.buildResult = this.getResult("BUILDRESULT"); this.buildResult == 0 {
		this.runResult = this.getLog("RUNRESULT")
	} else {
		this.runResult = "EC"
	}

	return map[string]interface{}{
		"build_log":    this.buildLog,
		"build_result": this.buildResult,
		"run_result":   this.runResult,
	}
}

func (this *Info) getLog(file string) string {
	path := filepath.Join(this.buildPath, this.sid, fmt.Sprintf("%d", this.id), file)
	if com.PathExist(path) {
		content, _ := com.ReadFile(path)
		return content
	} else {
		return ""
	}
}

// get the result
func (this *Info) getResult(file string) int {
	path := filepath.Join(this.buildPath, this.sid, fmt.Sprintf("%d", this.id), file)
	if com.PathExist(path) {
		content, _ := com.ReadFile(path)
		content = com.Strim(content)
		if result, err := strconv.Atoi(content); err != nil {
			log.Warnln(err)
			return -1
		} else {
			return result
		}
	} else {
		return -1
	}
}
