package main

import (
	"fmt"
)

func main() {
	cards := createNewDeck()

	hand, rest := deal(cards, 5)

	// cards.print()
	fmt.Println(hand)
	fmt.Println(len(rest))
	fmt.Println(len(cards))

	saveToFile(cards, "my_cards")
}
