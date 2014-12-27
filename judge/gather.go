package judge

import (
	"fmt"
	"github.com/gogather/com"
	"path/filepath"
	"strconv"
)

type Info struct {
	sid         string
	id          int
	buildLog    string
	buildResult bool
	runLog      string
	runResult   bool
	buildPath   string
}

func (this *Info) Gather(sid string, id int, buildPath string) map[string]interface{} {
	this.sid = sid
	this.id = id
	this.buildPath = buildPath

	this.buildLog = this.getLog("BUILD.LOG")
	if this.buildResult = this.getResult("BUILDRESULT"); this.buildResult {
		this.runLog = this.getLog("RUN.LOG")
		this.runResult = this.getResult("RUNRESULT")
	} else {
		this.runLog = ""
		this.runResult = false
	}

	return map[string]interface{}{
		"build_log":    this.buildLog,
		"build_result": this.buildResult,
		"run_log":      this.runLog,
		"run_result":   this.runResult,
	}
}

func (this *Info) getLog(file string) string {
	path := filepath.Join(this.buildPath, this.sid, fmt.Sprintf("%d", this.id), file)
	if com.PathExist(path) {
		return com.ReadFile(path)
	} else {
		return ""
	}
}

// get the result
func (this *Info) getResult(file string) bool {
	path := filepath.Join(this.buildPath, this.sid, fmt.Sprintf("%d", this.id), file)
	if com.PathExist(path) {
		content := com.ReadFile(path)
		if result, err := strconv.Atoi(content); err != nil {
			return false
		} else {
			if result == 0 {
				return true
			} else {
				return false
			}
		}
	} else {
		return false
	}
}
