package main

func main() {
	deck := newDeck()
	deck.print()
	deck.shuffle()
	deck.print()
}
