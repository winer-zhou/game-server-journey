package main

import "fmt"

// 注意返回值有名字：n
func f() (n int) {
	defer func() {
		n++ // 修改返回值
	}()

	return 1 // 实际上会先赋值给 n，然后 defer 把 n 变成 2
}

func main() {
	fmt.Println(f())
}
