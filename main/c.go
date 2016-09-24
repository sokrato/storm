package main

import (
	"fmt"
)

func main() {
	i := 0
	func() {
		i++
	}()
	fmt.Println(i)
}
