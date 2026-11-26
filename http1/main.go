package main

import (
	"fmt"
	"net/http"
)

func main() {
	// ── 注册路由：访问 "/" 时执行这个函数 ──
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "你好！你访问的是 %s\n", r.URL.Path)
		// w = ResponseWriter：往它里面写，就是写给浏览器的内容
		// r = *Request：请求信息（路径、方法、header 等）
	})

	// ── 注册路由：访问 "/health" 时执行 ──
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	fmt.Println("HTTP 服务器启动，监听 8081")

	// ── 启动（这行会阻塞，一直跑）──
	http.ListenAndServe(":8081", nil)
	// nil = 用默认的多路复用器（就是上面 HandleFunc 注册的那些）
}
