/*
Defines the core game state machine.

Notes:
- Players are generally referenced by a number.
*/
package game

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func Id() string {
	// TODO this should be a full length UUID, but 8 chars is easier to read in development.
	return uuid.NewString()[:8]
}

const MaxPlayers = 6

// An initial player action.
type Action struct {
	Type         ActionType `json:"type"`
	TargetPlayer string     `json:"target"`
}

// A block action, holding a card being used for a block and the player who initiated the block.
type Block struct {
	Card      Card   `json:"card"`
	Initiator string `json:"initiator"`
}

// A challenge action, holding the initiating player.
type Challenge struct {
	Initiator string `json:"initiator"`
}

// Represents the current state of a turn.
type TurnState string

const (
	Default             TurnState = "default"
	ActionPending       TurnState = "action_pending"
	BlockPending        TurnState = "block_pending"
	ExchangePending     TurnState = "exchange_pending"
	PlayerLostChallenge TurnState = "player_lost_challenge"
	LeaderLostChallenge TurnState = "leader_lost_challenge"
	PlayerKilled        TurnState = "player_killed"
	Finished            TurnState = "finished"
	PlayerWon           TurnState = "player_won"
)

// Represents the current game state.
type Game struct {
	Deck             []Card
	Players          map[string]*Player
	Winner           string
	Order            []string
	Leader           int
	TurnState        TurnState
	NextDeath        string
	PendingAction    Action
	PendingBlock     Block
	PendingChallenge Challenge
}

// Creates a new game with a shuffled deck.
func NewGame() Game {
	shuffled := ShuffleCards(Deck)

	game := Game{
		Deck:             shuffled,
		Players:          make(map[string]*Player),
		Winner:           "",
		Order:            []string{},
		Leader:           0,
		NextDeath:        "",
		PendingAction:    Action{},
		PendingBlock:     Block{},
		PendingChallenge: Challenge{},
		TurnState:        Default,
	}
	return game
}

// Adds a player to the game, returning their player number.
func (g *Game) AddPlayer(id string, name string) error {
	if len(g.Players) >= MaxPlayers {
		return errors.New("attempted to add player to a full game")
	}

	player := NewPlayer(id, name)
	g.Players[id] = &player
	g.Order = append(g.Order, id)
	return nil
}

func (g *Game) GetLeader() *Player {
	id := g.Order[g.Leader]
	return g.Players[id]
}

func (g *Game) GetPlayerByIndex(index int) (*Player, error) {
	if index > len(g.Players) {
		return nil, errors.New("attempted to access out of range player")
	}
	id := g.Order[index]
	player, ok := g.Players[id]
	if !ok {
		return nil, fmt.Errorf("player with id %s not found", id)
	}
	return player, nil
}

// Deals two cards to each players in the current game.
func (g *Game) Deal() {
	playerCount := len(g.Players)
	for i := range playerCount * 2 {
		// Remove the last card from the active deck and give it to the player.
		index := len(g.Deck) - 1
		card := g.Deck[index]
		g.Deck = g.Deck[:index]
		p, err := g.GetPlayerByIndex(i % playerCount)
		if err != nil {
			panic("player order not set up correctly")
		}
		p.GiveCard(card)
	}
}

// Transition from the default game state to ActionPending.
func (g *Game) AttemptAction(action Action) error {
	if !g.stateIn(Default) {
		return errors.New("action already in play")
	}
	leader := g.GetLeader()

	_, ok := g.Players[action.TargetPlayer]
	if action.TargetPlayer != "" && !ok {
		return errors.New("target player does not exist")
	}

	// Cost is always applied, even if an action is blocked or challenged.
	err := leader.PayForAction(action.Type)
	if err != nil {
		return err
	}

	g.PendingAction = action
	g.TurnState = ActionPending

	// Auto-commit income, as it isn't blockable.
	if g.PendingAction.Type == Income {
		g.CommitTurn()
	}
	return nil
}

// Attempt to block a pending action with an card.
func (g *Game) AttemptBlock(block Block) error {
	if !g.stateIn(ActionPending) {
		return errors.New("no action to block")
	}

	// Check the card being used blocks the current pending action.
	if !block.Card.BlocksAction(g.PendingAction.Type) {
		return errors.New("card does not block current pending action")
	}

	g.PendingBlock = block
	g.TurnState = BlockPending
	return nil
}

// If an action or block is pending, checks if the player who initiated the action has the correct card.
func (g *Game) Challenge(challenge Challenge) error {
	if !g.stateIn(ActionPending, BlockPending) {
		return errors.New("no action or block to challenge")
	}

	if _, ok := g.Players[challenge.Initiator]; !ok {
		return fmt.Errorf("invalid challenge initiator: %s", challenge.Initiator)
	}

	g.PendingChallenge = challenge
	leader := g.GetLeader()

	/*
		If an action is being challenged, check the leader has the necessary card for the
		current action.
	*/
	if g.TurnState == ActionPending {
		if leader.IsAllowedAction(g.PendingAction.Type) {
			g.TurnState = PlayerLostChallenge
			g.NextDeath = g.PendingChallenge.Initiator
			return nil
		}

		g.TurnState = LeaderLostChallenge
		g.NextDeath = leader.Id
		return nil
	}

	/*
		If a block is being challenged, check the blocker has the necessary card to block the
		current action.
	*/
	if g.TurnState == BlockPending {
		blocker := g.Players[g.PendingBlock.Initiator]
		if blocker.CanBlock(g.PendingAction.Type) {
			g.TurnState = LeaderLostChallenge
			g.NextDeath = leader.Id
			return nil
		}
		g.TurnState = PlayerLostChallenge
		g.NextDeath = g.PendingBlock.Initiator
		return nil
	}
	return nil
}

// If the game is in a state where a player must lose a card, this functions allows a card to be killed.
// Sets the card at index `card` to dead on the player who must die, depending on state.
func (g *Game) ResolveDeath(card int) error {
	if !g.stateIn(LeaderLostChallenge, PlayerLostChallenge, PlayerKilled) {
		return errors.New("no pending deaths to resolve")
	}

	if g.NextDeath == "" {
		return errors.New("id of next to die not set")
	}
	g.Players[g.NextDeath].KillCard(card)
	g.NextDeath = ""

	// If a player has lost a challenge, return to ActionPending
	switch g.TurnState {
	case PlayerLostChallenge:
		g.TurnState = ActionPending

	// If the leader has lost the challenge, the turn is over.
	case LeaderLostChallenge:
		g.TurnState = Finished

	// If a player has been killed (assassinated or coup'd), the turn is over.
	case PlayerKilled:
		g.TurnState = Finished
	}
	return nil
}

// Commits a turn, either confirming an action
func (g *Game) CommitTurn() error {
	if !g.stateIn(ActionPending, BlockPending) {
		return errors.New("no action or block is pending")
	}

	// If a block has not been challenged, end the turn.
	if g.TurnState == BlockPending {
		g.TurnState = Finished
		return nil
	}

	switch g.PendingAction.Type {
	case Income:
		g.GetLeader().AdjustCredits(1)
		g.TurnState = Finished

	case ForeignAid:
		g.GetLeader().AdjustCredits(2)
		g.TurnState = Finished

	case Revolt, Assassinate:
		g.NextDeath = g.PendingAction.TargetPlayer
		g.TurnState = PlayerKilled
	case Tax:
		g.GetLeader().AdjustCredits(3)
		g.TurnState = Finished

	// TODO not sure we even want this?
	case Exchange:
		// Draw two cards, select two from hand, return two to deck.
		g.TurnState = Finished

	case Steal:
		targetPlayer := g.PendingAction.TargetPlayer
		g.Players[targetPlayer].AdjustCredits(-2)
		g.GetLeader().AdjustCredits(2)
		g.TurnState = Finished
	default:
		return fmt.Errorf("tried to commit an unknown action %v", g.PendingAction.Type)
	}

	return nil
}

func (g *Game) EndTurn() error {
	if !g.stateIn(Finished) {
		return errors.New("turn not finished")
	}

	// Check for a winner.
	inPlay := []Player{}
	for _, player := range g.Players {
		if len(player.GetLivingCards()) != 0 {
			inPlay = append(inPlay, *player)
		}
	}

	if len(inPlay) == 1 {
		g.TurnState = PlayerWon
		g.Winner = inPlay[0].Id
		return nil
	}

	// Select the next player with living cards. Iterations are bounded to len(players) just in case.
	for range len(g.Players) {
		g.Leader = (g.Leader + 1) % len(g.Players)
		if len(g.GetLeader().GetLivingCards()) != 0 {
			break
		}
	}
	g.TurnState = Default
	g.PendingAction = Action{}
	g.PendingBlock = Block{}
	g.PendingChallenge = Challenge{}
	return nil
}

// Utility function to check if the game is in any of the passed states.
func (g *Game) stateIn(states ...TurnState) bool {
	for _, state := range states {
		if g.TurnState == state {
			return true
		}
	}
	return false
}
