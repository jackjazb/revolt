package main

import (
	"reflect"
	"testing"
)

func TestShuffle(t *testing.T) {
	t.Run("should shuffle a passed set of cards", func(t *testing.T) {
		initial := Deck
		shuffled := ShuffleCards(initial)
		if reflect.DeepEqual(initial, shuffled) {
			t.Errorf("cards not shuffled")
		}
	})
}
