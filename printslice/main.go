package main

import "fmt"

func main() {
	elements := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for _, element := range elements {
		if element%2 == 0 {
			fmt.Println("Even")
		} else {
			fmt.Println("Odd")
		}
	}
}
