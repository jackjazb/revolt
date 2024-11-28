package main

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	game := NewGame()

	t.Run("should set up a game in the default state", func(t *testing.T) {
		state := game.TurnState
		if state != Default {
			t.Error("incorrect game state, expected Default got", state)
		}
	})
}

func TestAddPlayer(t *testing.T) {
	game := NewGame()

	t.Run("should add players to the game", func(t *testing.T) {
		game.AddPlayer()
		game.AddPlayer()

		players := len(game.Players)
		if players != 2 {
			t.Error("incorrect player count, expected 2 got", players)
		}
	})

	t.Run("should reject too many players", func(t *testing.T) {
		game = NewGame()
		game.AddPlayer()
		game.AddPlayer()
		game.AddPlayer()
		game.AddPlayer()
		game.AddPlayer()
		game.AddPlayer()
		_, err := game.AddPlayer()

		if err == nil {
			t.Error("expected error adding player , got nil")
		}
	})
}

func TestDeal(t *testing.T) {
	game := NewGame()
	game.AddPlayer()
	game.AddPlayer()
	game.Deal()

	t.Run("should give each player two cards", func(t *testing.T) {
		cards := len(game.Players[0].Cards)
		if cards != 2 {
			t.Error("incorrect player card count, expected 2 got", cards)
		}
	})

	t.Run("should keep the right number of cards in play", func(t *testing.T) {
		playerOneCards := len(game.Players[0].Cards)
		playerTwoCards := len(game.Players[1].Cards)
		deck := len(game.Deck)
		cards := playerOneCards + playerTwoCards + deck
		if cards != len(Deck) {
			t.Errorf("wrong number of cards in play, expected %d got %d", len(Deck), cards)
		}
	})
}

func AttemptAction(t *testing.T) {
	game := NewGame()
	action := Action{
		Type:         Assassinate,
		TargetPlayer: 0,
	}
	game.AttemptAction(action)

	t.Run("should transition the game to the ActionPending state", func(t *testing.T) {
		if game.TurnState != ActionPending {
			t.Error("expected game to be in ActionPending, got", game.TurnState)
		}
	})
}
