package main

import "fmt"

func main() {
	x := 1

	defer fmt.Println("defer 里的 x =", x) // ← 注意这行
	//                                  ↑
	//                    参数 x 在【这行】就被复制成 1 了

	x = 100

	fmt.Println("后来的 x =", x)
}
