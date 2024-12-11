package main

import (
	"reflect"
	"revolt/game"
	"testing"
)

func TestToClientStateBroadCast(t *testing.T) {
	setup := func() GameInstance {
		i := NewGameInstance("0")
		i.Game.AddPlayer("0", "Player One")
		i.Game.AddPlayer("1", "Player Two")
		i.Game.Players["0"].Cards = append(
			i.Game.Players["0"].Cards,
			game.CardState{Card: game.Contessa, Alive: true},
			game.CardState{Card: game.Captain, Alive: false},
		)

		return i
	}

	t.Run("should return a client specific state broadcast", func(t *testing.T) {
		i := setup()

		broadcast := i.ToClientStateBroadcast(&Client{Id: "0"})

		if broadcast.Self.Id != "0" {
			t.Errorf("expected self.id to be 0, got: %s", broadcast.Self.Id)
		}
	})

	t.Run("should identify the current leader", func(t *testing.T) {
		i := setup()

		broadcast := i.ToClientStateBroadcast(&Client{Id: "0"})

		if !broadcast.Self.Leading {
			t.Errorf("expected self.leading to be true")
		}

		if broadcast.Peers[0].Leading {
			t.Errorf("expected peer to not be leading")
		}
	})

	t.Run("should include all a client's cards", func(t *testing.T) {
		i := setup()

		broadcast := i.ToClientStateBroadcast(&Client{Id: "0"})

		expectedCards := []game.CardState{
			{Card: game.Contessa, Alive: true},
			{Card: game.Captain, Alive: false},
		}

		if !reflect.DeepEqual(expectedCards, broadcast.Self.Cards) {
			t.Errorf("expected self.cards to be %v, got: %v", expectedCards, broadcast.Self.Cards)

		}
	})

	t.Run("should conceal living cards", func(t *testing.T) {
		i := setup()

		broadcast := i.ToClientStateBroadcast(&Client{Id: ""})

		expectedCards := []game.CardState{
			{Card: game.Captain, Alive: false},
		}

		if !reflect.DeepEqual(expectedCards, broadcast.Peers[0].Cards) {
			t.Errorf("expected self.cards to be %v, got: %v", expectedCards, broadcast.Self.Cards)

		}
	})
}
