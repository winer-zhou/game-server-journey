package main

import "fmt"

func main() {
	// 创建一个无缓冲的 channel，类型是 string
	ch := make(chan string)

	// 启动一个后台 goroutine，向 channel 发送数据
	go func() {
		ch <- "我是后台 goroutine" // 发送数据到 channel
	}()

	// 从 channel 接收数据并打印
	msg := <-ch
	fmt.Println("收到", msg)
}
