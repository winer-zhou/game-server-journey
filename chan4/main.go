package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// 启动一个 goroutine，向 ch1 发送数据
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "一号来了"
	}()

	// 启动另一个 goroutine，向 ch2 发送数据
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "二号来了"
	}()

	//等两次
	for i := 0; i < 2; i++ {
		select {
		case ch <- msg:
			fmt.Println("塞进去了")
		default:
			fmt.Println("队列满了,塞不进去了 踢掉这个客户端")
		case msg1 := <-ch1:
			fmt.Println("收到来自 ch1 的消息:", msg1)
		case msg2 := <-ch2:
			fmt.Println("收到来自 ch2 的消息:", msg2)
		}
	}
}
