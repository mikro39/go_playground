package main

import (
	"fmt"
)

func main() {
	var a [5]int //big old way to write arrays
	a[2] = 7     // setting an array postion 2 to int 7

	b := [5]int{4, 5, 1, 2, 3} // shorthand way to do the above array set

	c := []int{2, 6, 7, 4, 9} // open array ready for slice
	c = append(c, 7)          // appending an array to add the number "7"

	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)

	e := make(map[string]int) // building the map
	e["make"] = 0
	e["color"] = 1
	e["time"] = 2

	fmt.Println(e) // printing the map
	fmt.Println(e["time"]) // printing the array position within the map to "time"

	delete(e, "make") //showing that the "delete" command works to remove an item from the map
	fmt.Println(e)
}
