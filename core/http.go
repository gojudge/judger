package core

import (
	"github.com/gogather/com"
	"github.com/gogather/com/log"
	"net/http"
)

func HttpStart() {

	http.HandleFunc("/", HandleJsonRpc)

	err := http.ListenAndServe(":1005", nil)
	if err != nil {
		log.Warnln("ListenAndServe: ", err)
	} else {
		log.Blueln("Http Server Started!")
	}
}

func HandleJsonRpc(w http.ResponseWriter, r *http.Request) {

	p := make([]byte, r.ContentLength)
	r.Body.Read(p)

	content := string(p)

	log.Blueln(content)

	json, err := com.JsonDecode(content)

	if err != nil {
		log.Warnln("not json-rpc format")
		return
	}

	data := json.(map[string]interface{})

	ctrl, exists := RouterMap[data["action"].(string)]
	if !exists {
		log.Warnln("not exist")
		return
	}
	ctrl.Http(w, r)
}
