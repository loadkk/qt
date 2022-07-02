package test

import (
	"fmt"
	"net/http"
	"time"
)

// HttpServer 网页访问 localhost:9999 就能看到 Let's Go
func HttpServer() {
	callback := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Let's Go! %s", time.Now())
	}

	http.HandleFunc("/", callback)
	_ = http.ListenAndServe(":8888", nil) // 监听端口
}
