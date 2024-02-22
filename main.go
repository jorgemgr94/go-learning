package main

import (
	"fmt"
)

func main() {
	deck1 := createNewDeck()
	deck2 := newDeckFromFile("my_cards")

	hand, rest := deal(deck1, 5)

	fmt.Println(len(deck1))
	fmt.Println(len(hand))
	fmt.Println(len(rest))
	fmt.Println(len(deck2))

}
