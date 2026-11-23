package main
import "fmt"
func main(){
	n := []int{1,2,3}
	fmt.Println(n)

	n = append(n,4)
	fmt.Println(n)
	fmt.Println(n,len(n))
}