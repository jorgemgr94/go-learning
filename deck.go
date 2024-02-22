package main

import (
	"fmt"
	"os"
	"strings"
)

type deck []string

// == Deck functions =================================
func createNewDeck() deck {
	// array vs slice
	// array: fixed length list of things
	// slice: an array that can grow or shrink
	cards := deck{}
	cardSuits := []string{"Spades", "Diamonds", "Hearts", "Clubs"}
	cardValues := []string{"Ace", "Two", "Three", "Four"}

	for _, suit := range cardSuits {
		for _, value := range cardValues {
			cards = append(cards, value+" of "+suit)
		}
	}

	return cards
}

func saveToFile(d deck, filename string) {
	// byte slice is the representation of the string
	// as a ASCII character code (byte)
	os.WriteFile(filename, []byte(d.toString()), 0666)
}

// == Receiver functions =============================
/*
NOTE: this function is a receiver (d deck)
We're extending any variable of type "deck"
with the "print" method.
*/
func (d deck) print() {
	for i, card := range d {
		fmt.Println(i, card)
	}
}

func deal(d deck, handSize int) (deck, deck) {
	return d[:handSize], d[handSize:]
}

func (d deck) toString() string {
	return strings.Join(d, ",")
}
