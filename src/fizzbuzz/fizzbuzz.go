package main

import (
	"fmt"
)

func main() {
	for x := 1; x <= 120; x++ {
		if x%3 == 0 {
			fmt.Printf("fizz")
		}
		if x%5 == 0 {
			fmt.Printf("buzz")
		}
		if x%5 != 0 && x%3 != 0 {
			fmt.Printf("%d", x)
		}

		fmt.Printf("\n")
	}
}
