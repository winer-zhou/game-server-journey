package main

import "fmt"

func main() {
	// 创建一个无缓冲的 channel，类型是 string
	ch := make(chan int, 1)

	select {
	case ch <- 1:
		fmt.Println("塞进去了")
	default:
		fmt.Println("满了，塞不进去了")
	}

	select {
	case ch <- 2:
		fmt.Println("塞进去了")
	default:
		fmt.Println("满了，塞不进去了")
	}

	fmt.Println("取出来", <-ch)
}
