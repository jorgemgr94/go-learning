package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type deck struct {
	cards []string
}

// == Deck functions =================================
func newDeck() deck {
	// array vs slice
	// array: fixed length list of things
	// slice: an array that can grow or shrink
	deck := deck{cards: []string{}}
	cardSuits := []string{"Spades", "Diamonds", "Hearts", "Clubs"}
	cardValues := []string{"Ace", "Two", "Three", "Four"}

	for _, suit := range cardSuits {
		for _, value := range cardValues {
			deck.cards = append(deck.cards, value+" of "+suit)
		}
	}

	return deck
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
	return deck{cards: cards}
}

func deal(d deck, handSize int) (deck, deck) {
	return deck{cards: d.cards[:handSize]}, deck{cards: d.cards[handSize:]}
}

// == Receiver functions =============================
/*
NOTE: this function is a receiver (d deck)
We're extending any variable of type "deck"
with the "print" method.
*/
func (d deck) print() {
	fmt.Println("Printing deck:", d)
	for i, card := range d.cards {
		fmt.Println(i, card)
	}
}

func (d deck) toString() string {
	return strings.Join(d.cards, ",")
}

// Receiver functions
func (d deck) saveToFile(filename string) {
	// byte slice is the representation of the string
	// as a ASCII character code (byte)
	os.WriteFile(filename, []byte(d.toString()), 0666)
}

func (d deck) shuffle() {
	// NOTE: creating a seed for the random number generator
	//       using the current time as the seed
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for i := range d.cards {
		// generate a random number between 0 and len(d)
		newPosition := r.Intn(len(d.cards) - 1)

		// swap the current card with the card at the newPosition
		d.cards[i], d.cards[newPosition] = d.cards[newPosition], d.cards[i]
	}
}
