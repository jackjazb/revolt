package main

import (
	"reflect"
	"testing"
)

func TestGetCards(t *testing.T) {
	t.Run("should return a user's cards", func(t *testing.T) {
		player := Player{
			Cards: [2]CardState{
				{
					Card:  Ambassador,
					Alive: true,
				},
				{
					Card:  Captain,
					Alive: true,
				},
			},
		}
		expected := []Card{Ambassador, Captain}
		result := player.GetCards()

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %s expected %s", result, expected)
		}
	})

	t.Run("should ignore non living cards", func(t *testing.T) {
		player := Player{
			Cards: [2]CardState{
				{
					Card:  Duke,
					Alive: true,
				},
				{
					Card:  Captain,
					Alive: false,
				},
			},
		}
		expected := []Card{Duke}
		result := player.GetCards()

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %s expected %s", result, expected)
		}
	})
}
