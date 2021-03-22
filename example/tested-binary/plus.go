package main

import "fmt"

func main() {
	a := 5
	b := 7
	fmt.Printf("%d + %d = %d \n", a, b, plus(a, b))
}

func plus(a, b int) int {
	fmt.Printf("execute plus operation on arguments: %d, %d \n", a, b)
	return a + b
}
