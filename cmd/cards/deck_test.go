package main

import (
	"os"
	"testing"
)

func TestNewDeck(t *testing.T) {
	d := newDeck()

	if len(d.cards) != 16 {
		t.Errorf("Expected deck length of 16, but got %v", len(d.cards))
	}

	if d.cards[0] != "Ace of Spades" {
		t.Errorf("Expected first card of Ace of Spades, but got %v", d.cards[0])
	}

	if d.cards[len(d.cards)-1] != "Four of Clubs" {
		t.Errorf("Expected last card of Four of Clubs, but got %v", d.cards[len(d.cards)-1])
	}
}

func TestSaveToDeckAndNewDeckFromFile(t *testing.T) {
	// NOTE: remove any file with the name "_decktesting"
	os.Remove("_decktesting")

	deck := newDeck()
	deck.saveToFile("_decktesting")

	loadedDeck := newDeckFromFile("_decktesting")

	if len(loadedDeck.cards) != 16 {
		t.Errorf("Expected 16 cards in deck, got %v", len(loadedDeck.cards))
	}

	os.Remove("_decktesting")
}
