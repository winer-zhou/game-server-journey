package main

import "fmt"

func bubblesort(arr []int) {
	n := len(arr)

	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}

		}
	}
}
func main() {
	nums := []int{5, 2, 9, 1, 5, 6}
	bubblesort(nums)
	fmt.Println(nums)
}
