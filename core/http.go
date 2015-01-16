package core

import (
	"github.com/gogather/com/log"
	"net/http"
)

func HttpStart() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	err := http.ListenAndServe(":1005", nil)
	if err != nil {
		log.Warnln("ListenAndServe: ", err)
	} else {
		log.Blueln("Http Server Started!")
	}
}
