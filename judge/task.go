package judge

import (
	"errors"
	// "github.com/gogather/com/log"
	"fmt"
	"github.com/duguying/judger/core"
	"html"
	"runtime"
)

func AddTask(data map[string]interface{}) error {
	compiler := Compile{}
	var ok bool

	// HTML反转义
	code, ok := data["code"].(string)
	code = html.UnescapeString(code)
	if !ok {
		return errors.New("invalid code, should be string")
	}

	// language
	lang, ok := data["language"].(string)
	if !ok {
		return errors.New("invalid language, should be string")
	}

	// id
	id, ok := data["id"].(float64)
	if !ok {
		return errors.New("invalid language, should be string")
	}

	// session id
	sid, ok := data["sid"].(string)
	if !ok {
		return errors.New("invalid language, should be string")
	}

	// run compiling
	compiler.NewCompile()
	runPath, err := compiler.Run(code, lang, int(id), sid)

	if err == nil {
		// execute the binary in sandbox
		err = RunNativeInSandbox(core.C.Get(runtime.GOOS, "run_script"), fmt.Sprintf("%d", int(id)), runPath, 0, 0)
	}

	// TODO gather result

	return err
}
