package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("开始")

	ch := make(chan string)

	go func() {
		fmt.Println("后台 goroutine 准备发送数据")
		ch <- "数据"
		fmt.Println("后台 goroutine 发送数据完成")
	}()

	fmt.Println("主 goroutine 等待接收数据")
	time.Sleep(3 * time.Second) // 模拟主 goroutine 做其他事情

	fmt.Println("主 goroutine 准备接收数据")
	msg := <-ch
	fmt.Println("主 goroutine 收到数据:", msg)
}
