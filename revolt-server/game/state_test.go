package game

import (
	"fmt"
	"testing"
)

func TestNewGame(t *testing.T) {
	t.Run("should set up a game in the default state", func(t *testing.T) {
		g := NewGame()
		state := g.TurnState
		if state != Default {
			t.Error("incorrect game state, expected Default got", state)
		}
	})
}

func TestAddPlayer(t *testing.T) {
	t.Run("should add players to the game", func(t *testing.T) {
		g := NewGame()

		g.AddPlayer("0", "Test")
		g.AddPlayer("1", "Test")

		players := len(g.Players)
		if players != 2 {
			t.Error("incorrect player count, expected 2 got", players)
		}
	})

	t.Run("should reject too many players", func(t *testing.T) {
		g := NewGame()

		g.AddPlayer("0", "Test")
		g.AddPlayer("1", "Test")
		g.AddPlayer("2", "Test")
		g.AddPlayer("3", "Test")
		g.AddPlayer("4", "Test")
		g.AddPlayer("5", "Test")
		err := g.AddPlayer("6", "Test")

		if err == nil {
			t.Error("expected error adding player , got nil")
		}
	})
}

func TestGetPlayerById(t *testing.T) {
	t.Run("should return the player at a given point in the order", func(t *testing.T) {
		g := NewGame()

		g.AddPlayer("abc", "Test")
		g.AddPlayer("def", "Test")

		p, err := g.GetPlayerByIndex(0)
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		if p.Id != "abc" {
			t.Errorf("expected ID to be abc, got %s", p.Id)
		}
	})
}

func TestDeal(t *testing.T) {
	t.Run("should give each player two cards", func(t *testing.T) {
		g := NewGame()
		g.AddPlayer("0", "Test")
		g.AddPlayer("1", "Test")
		g.Deal()
		cards := len(g.Players["0"].Cards)
		if cards != 2 {
			t.Error("incorrect player card count, expected 2 got", cards)
		}
	})

	t.Run("should keep the right number of cards in play", func(t *testing.T) {
		g := NewGame()
		g.AddPlayer("0", "Test")
		g.AddPlayer("1", "Test")
		g.Deal()
		playerOneCards := len(g.Players["0"].Cards)
		playerTwoCards := len(g.Players["1"].Cards)
		deck := len(g.Deck)
		cards := playerOneCards + playerTwoCards + deck
		if cards != len(Deck) {
			t.Errorf("wrong number of cards in play, expected %d got %d", len(Deck), cards)
		}
	})
}

func TestAttemptAction(t *testing.T) {
	t.Run("should transition the game to the ActionPending state", func(t *testing.T) {
		g := NewGame()
		g.AddPlayer("0", "Test")
		g.Players["0"].AdjustCredits(1)

		err := g.AttemptAction(Action{
			Type:         Assassinate,
			TargetPlayer: "0",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		if g.TurnState != ActionPending {
			t.Error("expected game to be in ActionPending, got", g.TurnState)
		}
		if g.PendingAction.Type != Assassinate {
			t.Error("expected game to have assassination pending, got", g.PendingAction.Type)
		}

	})

	t.Run("should not allow targeting an out of range player", func(t *testing.T) {
		g := NewGame()
		g.AddPlayer("0", "Test")
		g.Players["0"].AdjustCredits(1)

		err := g.AttemptAction(Action{
			Type:         Assassinate,
			TargetPlayer: "no-one",
		})
		if err == nil {
			t.Error("expected an error, got nil")
		}
	})

	t.Run("should always apply action cost", func(t *testing.T) {
		g := NewGame()
		g.AddPlayer("0", "Test")
		g.Players["0"].AdjustCredits(1)

		err := g.AttemptAction(Action{
			Type:         Assassinate,
			TargetPlayer: "0",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		credits := g.Players["0"].Credits
		if credits != 0 {
			t.Errorf("expected 0 credits on player 1, got: %d", credits)
		}
	})

	t.Run("should not allow unaffordable actions", func(t *testing.T) {
		g := NewGame()
		g.AddPlayer("0", "Test")

		err := g.AttemptAction(Action{
			Type:         Revolt,
			TargetPlayer: "0",
		})
		if err == nil {
			t.Error("expected an error, got nil")
		}

		credits := g.Players["0"].Credits
		if credits != 2 {
			t.Errorf("expected 2 credits on player 1, got: %d", credits)
		}
	})
}

func TestAttemptBlock(t *testing.T) {
	setup := func() (Game, error) {
		g := NewGame()
		g.AddPlayer("0", "Test")
		g.AddPlayer("1", "Test")

		g.Players["0"].AdjustCredits(1)

		err := g.AttemptAction(Action{
			Type:         Assassinate,
			TargetPlayer: "1",
		})
		if err != nil {
			return Game{}, err
		}
		return g, nil
	}

	t.Run("should transition the game to the BlockPending state", func(t *testing.T) {
		g, err := setup()
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		err = g.AttemptBlock(Block{
			Card:      Contessa,
			Initiator: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		if g.TurnState != BlockPending {
			t.Errorf("expected game to be in BlockPending, got: %s", g.TurnState)
		}
	})

	t.Run("should set the pending block", func(t *testing.T) {
		g, err := setup()
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		err = g.AttemptBlock(Block{
			Card:      Contessa,
			Initiator: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		if g.PendingBlock.Card != Contessa {
			t.Errorf("expected block card to be Contessa, go: %s", g.PendingBlock.Card)
		}
		if g.PendingBlock.Initiator != "1" {
			t.Errorf("expected initiator to be 1, got: %s", g.PendingBlock.Initiator)
		}
	})
	t.Run("should not allow cards which do not block the pending action to be used", func(t *testing.T) {
		g, err := setup()
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		err = g.AttemptBlock(Block{
			Card:      Captain,
			Initiator: "1",
		})
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if g.TurnState != ActionPending {
			t.Errorf("expected game to still be in ActionPending, got: %s", g.TurnState)
		}
	})

	t.Run("should not allow blocking unblockable actions", func(t *testing.T) {
		g := NewGame()
		g.AddPlayer("0", "Test")
		g.AddPlayer("1", "Test")

		g.Players["0"].AdjustCredits(5)

		err := g.AttemptAction(Action{
			Type:         Revolt,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		err = g.AttemptBlock(Block{
			Card:      Contessa,
			Initiator: "1",
		})
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if g.TurnState != ActionPending {
			t.Errorf("expected game to still be in ActionPending, got: %s", g.TurnState)
		}
	})

	t.Run("should fail if the game is not in ActionPending", func(t *testing.T) {
		g := NewGame()
		err := g.AttemptBlock(Block{
			Card:      Captain,
			Initiator: "1",
		})
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestChallenge(t *testing.T) {
	setup := func() Game {
		g := NewGame()
		g.AddPlayer("0", "Test")
		g.AddPlayer("1", "Test")
		return g
	}

	t.Run("should transition to PlayerLostChallenge if the leader is allowed their action", func(t *testing.T) {
		g := setup()
		g.Players["0"].Cards = append(g.Players["0"].Cards, CardState{
			Card:  Captain,
			Alive: true,
		})

		err := g.AttemptAction(Action{
			Type:         Steal,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.Challenge(Challenge{
			Initiator: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		if g.TurnState != PlayerLostChallenge {
			t.Errorf("expected game to be in PlayerLostChallenge, got: %s", g.TurnState)
		}
	})

	t.Run("should transition to LeaderLostChallenge if the leader is not allowed their action", func(t *testing.T) {
		g := setup()

		err := g.AttemptAction(Action{
			Type:         Steal,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.Challenge(Challenge{
			Initiator: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		if g.TurnState != LeaderLostChallenge {
			t.Errorf("expected game to be in LeaderLostChallenge, got: %s", g.TurnState)
		}
	})

	t.Run("should transition to LeaderLostChallenge if the block initiator is allowed to block the leader", func(t *testing.T) {
		g := setup()
		g.Players["0"].AdjustCredits(1)
		g.Players["1"].Cards = append(g.Players["1"].Cards, CardState{
			Card:  Contessa,
			Alive: true,
		})

		err := g.AttemptAction(Action{
			Type:         Assassinate,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.AttemptBlock(Block{
			Card:      Contessa,
			Initiator: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.Challenge(Challenge{
			Initiator: "0",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		if g.TurnState != LeaderLostChallenge {
			t.Errorf("expected game to be in LeaderLostChallenge, got: %s", g.TurnState)
		}
	})

	t.Run("should transition to PlayerLostChallenge if the block initiator is not allowed to block the leader", func(t *testing.T) {
		g := setup()
		g.Players["0"].AdjustCredits(1)

		err := g.AttemptAction(Action{
			Type:         Assassinate,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.AttemptBlock(Block{
			Card:      Contessa,
			Initiator: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.Challenge(Challenge{
			Initiator: "0",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		if g.TurnState != PlayerLostChallenge {
			t.Errorf("expected game to be in PlayerLostChallenge, got: %s", g.TurnState)
		}
	})
}

func TestCommitTurn(t *testing.T) {
	setup := func() Game {
		g := NewGame()
		g.AddPlayer("0", "Test")
		g.AddPlayer("1", "Test")
		return g
	}

	t.Run("should transition from BlockPending to Finished if the block is not challenged", func(t *testing.T) {
		g := setup()

		err := g.AttemptAction(Action{
			Type:         Steal,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		err = g.AttemptBlock(Block{
			Card:      Captain,
			Initiator: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.CommitTurn()
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		if g.TurnState != Finished {
			t.Errorf("expected game to be in Finished, got: %s", g.TurnState)
		}
	})

	t.Run("should apply Income", func(t *testing.T) {
		g := setup()

		err := g.AttemptAction(Action{
			Type: Income,
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.CommitTurn()
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		credits := g.Players["0"].Credits
		if credits != 3 {
			t.Errorf("expected player to have 3 credits, got: %d", credits)
		}
	})

	t.Run("should apply Foreign Aid", func(t *testing.T) {
		g := setup()

		err := g.AttemptAction(Action{
			Type: ForeignAid,
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.CommitTurn()
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		credits := g.Players["0"].Credits
		if credits != 4 {
			t.Errorf("expected player to have 4 credits, got: %d", credits)
		}
	})

	t.Run("should apply Tax", func(t *testing.T) {
		g := setup()

		err := g.AttemptAction(Action{
			Type: Tax,
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.CommitTurn()
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		credits := g.Players["0"].Credits
		if credits != 5 {
			t.Errorf("expected player to have 5 credits, got: %d", credits)
		}
	})

	t.Run("should apply assassinations", func(t *testing.T) {
		g := setup()
		g.Players["0"].AdjustCredits(1)

		err := g.AttemptAction(Action{
			Type:         Assassinate,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.CommitTurn()
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		if g.NextDeath != "1" {
			t.Errorf("expected next death to be in 1, got: %s", g.NextDeath)
		}

		if g.TurnState != PlayerKilled {
			t.Errorf("expected game to be in PlayerKilled, got: %s", g.TurnState)
		}
	})

	t.Run("should apply revolts", func(t *testing.T) {
		g := setup()
		g.Players["0"].AdjustCredits(5)

		err := g.AttemptAction(Action{
			Type:         Revolt,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.CommitTurn()
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		if g.TurnState != PlayerKilled {
			t.Errorf("expected game to be in PlayerKilled, got: %s", g.TurnState)
		}
	})

	t.Run("should apply theft", func(t *testing.T) {
		g := setup()

		err := g.AttemptAction(Action{
			Type:         Steal,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.CommitTurn()
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		p1credits := g.Players["0"].Credits
		if p1credits != 4 {
			t.Errorf("expected player 1 to have 4 credits, got: %d", p1credits)
		}
		p2credits := g.Players["1"].Credits
		if p2credits != 0 {
			t.Errorf("expected player 2 to have 4 credits, got: %d", p2credits)
		}
	})
}

func TestResolveDeath(t *testing.T) {
	setup := func() Game {
		g := NewGame()
		g.AddPlayer("0", "Test")
		g.AddPlayer("1", "Test")
		return g
	}

	t.Run("should resolve a pending death caused by a successful action", func(t *testing.T) {
		g := setup()
		g.Players["0"].AdjustCredits(5)
		g.Players["1"].Cards = append(g.Players["1"].Cards, CardState{
			Card:  Contessa,
			Alive: true,
		})

		err := g.AttemptAction(Action{
			Type:         Revolt,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.CommitTurn()
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		err = g.ResolveDeath(0)
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		if g.Players["1"].Cards[0].Alive {
			t.Errorf("expected card to be dead")
		}
	})

	t.Run("should resolve a pending death caused by a successfully challenged block", func(t *testing.T) {
		g := setup()
		g.Players["0"].Cards = append(g.Players["0"].Cards, CardState{
			Card:  Assassin,
			Alive: true,
		})
		g.Players["0"].AdjustCredits(1)

		g.Players["1"].Cards = append(g.Players["1"].Cards, CardState{
			Card:  Captain,
			Alive: true,
		})

		err := g.AttemptAction(Action{
			Type:         Assassinate,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.AttemptBlock(Block{
			Card:      Contessa,
			Initiator: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.Challenge(Challenge{
			Initiator: "0",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		fmt.Println(g.TurnState)
		err = g.ResolveDeath(0)
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		if g.Players["1"].Cards[0].Alive {
			t.Errorf("expected card to be dead")
		}
	})

	t.Run("should resolve a pending death caused by a successfully challenged action", func(t *testing.T) {
		g := setup()
		g.Players["0"].Cards = append(g.Players["0"].Cards, CardState{
			Card:  Contessa,
			Alive: true,
		})
		g.Players["0"].AdjustCredits(1)

		err := g.AttemptAction(Action{
			Type:         Assassinate,
			TargetPlayer: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}

		err = g.Challenge(Challenge{
			Initiator: "1",
		})
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		err = g.ResolveDeath(0)
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		if g.Players["0"].Cards[0].Alive {
			t.Errorf("expected card to be dead")
		}
	})
}

func TestEndTurn(t *testing.T) {
	setup := func() Game {
		g := NewGame()
		g.AddPlayer("0", "Test")
		g.AddPlayer("1", "Test")
		g.AddPlayer("2", "Test")
		return g
	}
	t.Run("should advance the leader and transition the game back to the default state", func(t *testing.T) {
		g := setup()
		err := g.AttemptAction(Action{Type: Income})
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		err = g.CommitTurn()
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		err = g.EndTurn()
		if err != nil {
			t.Errorf("got error: %s", err)
		}
		if g.Leader != 1 {
			t.Errorf("expected 1 to be leader, go %d", g.Leader)
		}
	})

	t.Run("should return the leader to 0 after a sufficient number of turns", func(t *testing.T) {
		g := setup()
		for range 3 {
			err := g.AttemptAction(Action{Type: Income})
			if err != nil {
				t.Errorf("got error: %s", err)
			}
			err = g.CommitTurn()
			if err != nil {
				t.Errorf("got error: %s", err)
			}
			err = g.EndTurn()
			if err != nil {
				t.Errorf("got error: %s", err)
			}
		}

		if g.Leader != 0 {
			t.Errorf("expected 0 to be leader, go %d", g.Leader)
		}
	})
}
