package main

func main() {
	cards := newDeck()
	// hand, rest := deal(cards, 5)
	cards.shuffle()
	cards.print()
}
