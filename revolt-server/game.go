/*
Defines the core game logic.

Notes:
- Players are generally referenced by a number.
*/
package main

import (
	"errors"
	"fmt"
	"log"
	"slices"
)

const MaxPlayers = 6

// An initial player action.
type Action struct {
	Type         ActionType `json:"type"`
	TargetPlayer int        `json:"targetPlayer"`
}

// A block action, holding a card being used for a block and the player who initiated the block.
type Block struct {
	Card      Card `json:"card"`
	Initiator int  `json:"initiator"`
}

// A challenge action, holding the initiating player.
type Challenge struct {
	Initiator int `json:"initiator"`
}

// Represents the current state of a turn.
type TurnState int

const (
	Default TurnState = iota
	ActionPending
	BlockPending
	ExchangePending
	PlayerLostChallenge
	LeaderLostChallenge
	PlayerKilled
	Finished
)

// Formats a turn state.
func (i TurnState) String() string {
	return []string{
		"Default",
		"ActionPending",
		"BlockPending",
		"ExchangePending",
		"PlayerLostChallenge",
		"LeaderLostChallenge",
		"PlayerKilled",
		"Finished",
	}[i]
}

// Represents the current game state.
type Game struct {
	Deck             []Card
	Players          []Player
	Leader           int
	TurnState        TurnState
	PendingAction    Action
	PendingBlock     Block
	PendingChallenge Challenge
}

// Creates a new game with a shuffled deck.
func NewGame() Game {
	shuffled := ShuffleCards(Deck)

	game := Game{
		Deck:             shuffled,
		Players:          []Player{},
		Leader:           0,
		PendingAction:    Action{},
		PendingBlock:     Block{},
		PendingChallenge: Challenge{},
		TurnState:        Default,
	}
	return game
}

// Adds a player to the game, returning their player number.
func (g *Game) AddPlayer() (int, error) {
	if len(g.Players) >= MaxPlayers {
		return 0, errors.New("attempted to add player to a full game")
	}
	player := Player{
		Cards:   [2]CardState{},
		Credits: 2,
	}
	g.Players = append(g.Players, player)
	log.Println(g.Players)
	return len(g.Players) - 1, nil
}

// Deals two cards to each players in the current game.
func (g *Game) Deal() {
	playerCount := len(g.Players)
	fmt.Println(g.Players)
	for i := range playerCount * 2 {
		// Remove the last card from the active deck and give it to the player.
		index := len(g.Deck) - 1
		card := g.Deck[index]
		g.Deck = g.Deck[:index]

		cardState := CardState{
			Card:  card,
			Alive: true,
		}

		g.Players[i%playerCount].Cards[i/playerCount] = cardState
	}
	fmt.Printf("players after deal: %+v", g.Players)
}

// Transition from the default game state to ActionPending.
func (g *Game) AttemptAction(action Action) error {
	if !g.StateIn(Default) {
		return errors.New("action already in play")
	}
	leader := &g.Players[g.Leader]

	if action.TargetPlayer > len(g.Players)-1 {
		return errors.New("target player out of range")
	}

	// Cost is always applied, even if an action is blocked or challenged.
	if cost, ok := ActionCost[action.Type]; ok {
		if cost > leader.Credits {
			return errors.New("cannot afford action")
		}
		leader.Credits -= cost
	}

	g.PendingAction = action
	g.TurnState = ActionPending
	return nil
}

// Attempt to block a pending action with an card.
func (g *Game) AttemptBlock(block Block) error {
	if !g.StateIn(ActionPending) {
		return errors.New("no action to block")
	}

	// Check the card being used blocks the current pending action.
	blocks, ok := Blocks[block.Card]
	if !ok || blocks != g.PendingAction.Type {
		return errors.New("card does not block current pending action")
	}

	g.PendingBlock = block
	return nil
}

// If an action or block is pending, checks if the player who initiated the action has the correct card.
func (g *Game) Challenge(challenge Challenge) error {
	if !g.StateIn(ActionPending, BlockPending) {
		return errors.New("no action or block to challenge")
	}

	g.PendingChallenge = challenge

	/*
		If an action is being challenged, check the leader has the necessary card for the
		current action.
	*/
	if g.TurnState == ActionPending {
		leader := g.Players[g.Leader]
		requiredCard, ok := RequiredCard[g.PendingAction.Type]

		// If the current action is not in RequiredCards, no card is needed and the challenge fails.
		if !ok {
			g.TurnState = PlayerLostChallenge
			return nil
		}
		if slices.Contains(leader.GetCards(), requiredCard) {
			g.TurnState = PlayerLostChallenge
			return nil
		}
		g.TurnState = LeaderLostChallenge
		return nil
	}

	/*
		If a block is being challenged, check the blocker has the necessary card to block the
		current action.
	*/
	if g.TurnState == BlockPending {
		blocker := g.Players[g.PendingBlock.Initiator]
		requiredCards, ok := BlockedBy[g.PendingAction.Type]
		// If the action cannot be blocked, the challenge fails (shouldn't get here anyway)
		if !ok {
			g.TurnState = PlayerLostChallenge
			return nil
		}
		for _, card := range requiredCards {
			if slices.Contains(blocker.GetCards(), card) {
				g.TurnState = LeaderLostChallenge
				return nil
			}
		}
		g.TurnState = PlayerLostChallenge
		return nil
	}
	return nil
}

// If the game is in a state where a player must lose a card, this functions allows a card to be killed.
// Sets the card at index `card` to dead on the player who must die, depending on state.
func (g *Game) ResolveDeath(card int) error {
	if !g.StateIn(LeaderLostChallenge, PlayerLostChallenge, PlayerKilled) {
		return errors.New("no pending deaths to resolve")
	}

	// If a player has lost a challenge, return to ActionPending
	if g.TurnState == PlayerLostChallenge {
		g.Players[g.PendingChallenge.Initiator].Cards[card].Alive = false
		g.TurnState = ActionPending
		return nil
	}

	// If the leader has lost the challenge, the turn is over.
	if g.TurnState == LeaderLostChallenge {
		g.Players[g.Leader].Cards[card].Alive = false
		g.TurnState = Finished
		return nil
	}

	// If a player has been killed (assassinated or coup'd), the turn is over.
	if g.TurnState == PlayerKilled {
		g.Players[g.PendingAction.TargetPlayer].Cards[card].Alive = false
		g.TurnState = Finished
		return nil
	}
	return nil
}

// Commits a turn, either confirming an action
func (g *Game) CommitTurn() error {
	if !g.StateIn(ActionPending, BlockPending) {
		return errors.New("no action or block is pending")
	}

	// If a block has not been challenged, end the turn.
	if g.TurnState == BlockPending {
		g.TurnState = Finished
		return nil
	}

	switch g.PendingAction.Type {
	case Income:
		g.Players[g.Leader].Credits += 1
		g.TurnState = Finished

	case ForeignAid:
		g.Players[g.Leader].Credits += 2
		g.TurnState = Finished

	case Revolt, Assassinate:
		g.TurnState = PlayerKilled
	case Tax:
		g.Players[g.Leader].Credits += 3
		g.TurnState = Finished

	case Exchange:
		// Draw two cards, select two from hand, return two to deck.
		g.TurnState = Finished

	case Steal:
		targetPlayer := g.PendingAction.TargetPlayer
		g.Players[targetPlayer].Credits -= 2
		g.Players[g.Leader].Credits += 2
		g.TurnState = Finished
	default:
		return fmt.Errorf("tried to commit an unknown action %d", g.PendingAction.Type)
	}

	return nil
}

func (g *Game) EndTurn() error {
	if !g.StateIn(Finished) {
		return errors.New("turn not finished")
	}
	g.Leader = (g.Leader + 1) % len(g.Players)
	g.TurnState = Default
	return nil
}

// Utility function to check if the game is in any of the passed states.
func (g *Game) StateIn(states ...TurnState) bool {
	for _, state := range states {
		if g.TurnState == state {
			return true
		}
	}
	return false
}
