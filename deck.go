package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type deck []string

// == Deck functions =================================
func newDeck() deck {
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

func (d deck) saveToFile(filename string) {
	// byte slice is the representation of the string
	// as a ASCII character code (byte)
	os.WriteFile(filename, []byte(d.toString()), 0666)
}

func newDeckFromFile(filename string) deck {
	// byte slice is the representation of the string
	// as a ASCII character code (byte)
	bs, error := os.ReadFile(filename)
	if error != nil {
		fmt.Println("Error: ", error)
		os.Exit(1)
	}

	cards := strings.Split(string(bs), ",")
	return deck(cards)
}

func (d deck) shuffle() {
	// NOTE: creating a seed for the random number generator
	//       using the current time as the seed
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range d {
		// generate a random number between 0 and len(d)
		newPosition := r.Intn(len(d) - 1)

		// swap the current card with the card at the newPosition
		d[i], d[newPosition] = d[newPosition], d[i]
	}
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
