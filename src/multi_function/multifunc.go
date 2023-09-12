package main

import (
	"fmt"
)

func main() {
	fmt.Println("color")

	result := sum(2, 5) //values to pass to "sum" - in this ecample, as x, y
	fmt.Println(result)
}

// takes x int and y int and passes them into the results, in this case, the func "sum" returns the operation x + y
func sum(x int, y int) int {
	return x + y
}