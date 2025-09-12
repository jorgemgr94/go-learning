package main

import "fmt"

func maps() {
	fmt.Println("// == Maps ===========================================")
	// declaring a map with var
	// var colors2 map[string]string

	// declaring a map with make
	// colors3 := make(map[string]string)
	// colors3["red"] = "#ff0000"

	colors := map[string]string{
		"red":   "#ff0000",
		"green": "#4bf745",
		"white": "#ffffff",
	}

	fmt.Println(colors)

	// == manipulating a map ====================
	// deleting a key from a map
	delete(colors, "red")
	fmt.Println(colors)

	// == Iterating over a map ==================
	printMap(colors)
}

func printMap(c map[string]string) {
	// color => key, hex => value
	for color, hex := range c {
		fmt.Println("Hex code for", color, "is", hex)
	}
}
