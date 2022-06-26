package test

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type countHeader struct {
	mu sync.Mutex
	n  int
}

// ServeHTTP 这个函数不能改名字，必须严格遵守，这是Handle的要求，这块是个大坑
func (h *countHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//加解锁
	h.mu.Lock()
	defer h.mu.Unlock()
	h.n++
	_, _ = fmt.Fprintf(w, "count is %d\n", h.n)
}

// HttpServer2 打开网页localhost:9999，就能看到hello go！的文字了
func HttpServer2() {
	count := new(countHeader)
	fmt.Printf("count: %v\n", count)
	http.Handle("/hello", count)
	log.Fatal(http.ListenAndServe(":9999", nil))
}
