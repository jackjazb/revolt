package main

import "testing"

func TestNewGame(t *testing.T) {
	game, err := NewGame(2)
	if err != nil {
		t.Error("error creating game", err)
	}

	t.Run("should set up a game with the correct player count", func(t *testing.T) {
		players := len(game.Players)
		if players != 2 {
			t.Error("incorrect player count, expected 2 got", players)
		}
	})

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

	t.Run("should set up a game in the default state", func(t *testing.T) {
		state := game.TurnState
		if state != Default {
			t.Error("incorrect game state, expected Default got", state)
		}
	})

	t.Run("should reject too many players", func(t *testing.T) {
		game, err = NewGame(20)
		if err == nil {
			t.Error("expected error creating game, got nil")
		}
	})
}

func AttemptAction(t *testing.T) {
	game, err := NewGame(2)
	if err != nil {
		t.Error("error creating game", err)
	}

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
