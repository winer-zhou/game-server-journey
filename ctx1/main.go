package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 创建一个带有超时的 context，超时时间为 2 秒
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 确保在函数退出时释放资源

	done := make(chan string) // 创建一个 channel 用于接收任务完成的信号

	// 启动一个 goroutine，模拟一个耗时操作
	go func() {
		time.Sleep(5 * time.Second) // 模拟耗时操作
		done <- "任务完成"
	}()
	// 等待任务完成或超时
	select {
	case result := <-done:
		fmt.Println("成功", result)
	case <-ctx.Done():
		fmt.Println("操作超时:", ctx.Err())
	}
}
