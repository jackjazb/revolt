package game

import (
	"reflect"
	"testing"
)

func TestGetCards(t *testing.T) {
	t.Run("should return a user's cards", func(t *testing.T) {
		player := NewPlayer("id", "Test")
		player.GiveCard(Ambassador)
		player.GiveCard(Captain)

		expected := []Card{Ambassador, Captain}
		result := player.GetLivingCards()

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %s expected %s", result, expected)
		}
	})

	t.Run("should ignore non living cards", func(t *testing.T) {
		player := NewPlayer("id", "Test")
		player.GiveCard(Duke)
		player.GiveCard(Captain)
		player.KillCard(1)

		expected := []Card{Duke}
		result := player.GetLivingCards()

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %s expected %s", result, expected)
		}
	})
}

func TestGetAllowedActions(t *testing.T) {
	t.Run("should return a list of actions granted by a player's cards", func(t *testing.T) {
		player := NewPlayer("id", "Test")
		player.GiveCard(Captain)

		expected := []ActionType{Income, ForeignAid, Revolt, Steal}
		result := player.GetAllowedActions()

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %s expected %s", result, expected)
		}
	})

	t.Run("should return the default allowed actions if no other actions are allowed", func(t *testing.T) {
		player := NewPlayer("id", "Test")

		expected := []ActionType{Income, ForeignAid, Revolt}
		result := player.GetAllowedActions()

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %s expected %s", result, expected)
		}
	})

	t.Run("should ignore dead cards", func(t *testing.T) {
		player := NewPlayer("id", "Test")
		player.GiveCard(Captain)
		player.GiveCard(Assassin)
		player.KillCard(0)

		expected := []ActionType{Income, ForeignAid, Revolt, Assassinate}
		result := player.GetAllowedActions()

		if !reflect.DeepEqual(expected, result) {
			t.Errorf("got %s expected %s", result, expected)
		}
	})
}

func TestIsAllowedAction(t *testing.T) {
	t.Run("should test if a player is allowed to perform an action", func(t *testing.T) {
		player := NewPlayer("id", "Test")
		player.GiveCard(Captain)
		allowed := player.IsAllowedAction(Steal)
		if !allowed {
			t.Errorf("should have allowed action %s", Steal)
		}
	})

	t.Run("should allow default actions", func(t *testing.T) {
		player := NewPlayer("id", "Test")
		action := Income
		allowed := player.IsAllowedAction(action)
		if !allowed {
			t.Errorf("should have allowed action %s", action)
		}
	})
}

func TestPayForAction(t *testing.T) {
	t.Run("should deduct the credits for an action", func(t *testing.T) {
		player := NewPlayer("id", "Test")
		player.AdjustCredits(1)
		player.GiveCard(Assassin)
		player.PayForAction(Assassinate)

		if player.Credits != 0 {
			t.Errorf("expected 0 credits, got %d", player.Credits)
		}
	})

	t.Run("should error if the player cannot afford the action", func(t *testing.T) {
		player := NewPlayer("id", "Test")
		player.GiveCard(Assassin)
		err := player.PayForAction(Revolt)

		if player.Credits != 2 {
			t.Errorf("expected 2 credits, got %d", player.Credits)
		}
		if err == nil {
			t.Error("expected an error when paying for the action")
		}
	})
}

func TestCanBlock(t *testing.T) {
	t.Run("should return true if the player can block a given action", func(t *testing.T) {
		player := NewPlayer("id", "Test")
		player.GiveCard(Ambassador)
		action := Steal
		allowed := player.CanBlock(action)
		if !allowed {
			t.Errorf("expected to be able to block %s", action)
		}
	})

	t.Run("should return false if the player cannot block a given action", func(t *testing.T) {
		player := NewPlayer("id", "Test")
		player.GiveCard(Ambassador)
		action := Assassinate
		allowed := player.CanBlock(action)
		if allowed {
			t.Errorf("expected to be unable to block %s", action)
		}
	})
}
