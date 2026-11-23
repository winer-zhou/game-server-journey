package main

import "fmt"

// 两个数相加
func Add(a, b int) int {
	return a + b
}

// IsEven 判断是不是偶数
func IsEven(n int) bool {
	return n%2 == 0
}

func main() {
	fmt.Println(Add(1, 2))
	fmt.Println(IsEven(4))
}
