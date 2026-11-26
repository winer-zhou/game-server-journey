package main

import "fmt"

// ══════════════════════════════════════════════════════════
// ① 泛型函数：一份代码，支持多种类型
// ══════════════════════════════════════════════════════════
// [T int | float64 | string]
//
//	T           = 类型参数名（约定用 T，Type 的意思）
//	int|float64|string = 类型约束：只允许这三种
func Max[T int | float64 | string](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// ══════════════════════════════════════════════════════════
// ② 泛型结构体：一个"什么都能装"的栈
// ══════════════════════════════════════════════════════════
type Stack[T any] struct {
	//                    ↑ any = 任意类型（等同于 interface{}）
	items []T // 切片里装的类型是 T
}

// Push：入栈
func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop：出栈（返回值和"是否成功"）
func (s *Stack[T]) Pop() (T, bool) {
	var zero T // zero = T 类型的"零值"
	// int 的零值是 0，string 是 ""，指针是 nil

	if len(s.items) == 0 {
		return zero, false
	}

	v := s.items[len(s.items)-1]       // 取最后一个
	s.items = s.items[:len(s.items)-1] // 删掉最后一个
	return v, true
}

func main() {
	// ══════ 泛型函数：自动推断类型 ══════
	fmt.Println("Max int:     ", Max(3, 5))
	fmt.Println("Max float:   ", Max(3.14, 2.71))
	fmt.Println("Max string:  ", Max("abc", "abd"))

	// ══════ 泛型结构体：装 int ══════
	intStack := Stack[int]{}
	intStack.Push(10)
	intStack.Push(20)

	v, ok := intStack.Pop()
	fmt.Printf("int 栈弹出: %d (成功=%v)\n", v, ok)

	// ══════ 泛型结构体：装 string ══════
	strStack := Stack[string]{}
	strStack.Push("hello")
	strStack.Push("world")

	s, ok2 := strStack.Pop()
	fmt.Printf("string 栈弹出: %s (成功=%v)\n", s, ok2)
}
