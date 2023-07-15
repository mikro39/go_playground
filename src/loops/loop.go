package main

import (
	"fmt"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}

	fmt.Println("Now for a While loop") //The following sets up a while loop
	a := 0
	for a < 5 {
		fmt.Println(a)
		a++
	}

	fmt.Println("Now for looping an array")
	testarray := []string{"ab", "bb", "cb"}

	for index, value := range testarray {
		fmt.Println("index", index, "values", value)
	}
}
