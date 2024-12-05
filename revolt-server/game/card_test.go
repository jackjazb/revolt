package game

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

func TestBlocksAction(t *testing.T) {
	t.Run("should return true if a card blocks a given action", func(t *testing.T) {
		action := ForeignAid
		blocks := Duke.BlocksAction(action)
		if !blocks {
			t.Errorf("expected to block %s", action)
		}
	})

	t.Run("should return false if passed an un-blockable", func(t *testing.T) {
		action := Income
		blocks := Duke.BlocksAction(action)
		if blocks {
			t.Errorf("expected not to block %s", action)
		}
	})

	t.Run("should return false if a card does not block a given action", func(t *testing.T) {
		action := Assassinate
		blocks := Duke.BlocksAction(action)
		if blocks {
			t.Errorf("expected not to block %s", action)
		}
	})
}
