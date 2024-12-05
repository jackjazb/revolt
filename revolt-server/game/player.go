package game

import (
	"errors"
	"slices"
)

// Defines a single player.
type Player struct {
	Id      string      `json:"id"`
	Name    string      `json:"name"`
	Cards   []CardState `json:"cards"`
	Credits int         `json:"credits"`
}

func NewPlayer(id string, name string) Player {
	return Player{
		Id:      id,
		Name:    name,
		Cards:   []CardState{},
		Credits: 2,
	}
}

// Returns a player's cards.
func (p *Player) GetLivingCards() []Card {
	cards := []Card{}
	for _, card := range p.Cards {
		if card.Alive {
			cards = append(cards, card.Card)
		}
	}
	return cards
}

// Returns a player's dead cards as CardState structs.
func (p *Player) GetDeadCards() []CardState {
	cards := []CardState{}
	for _, card := range p.Cards {
		if !card.Alive {
			cards = append(cards, card)
		}
	}
	return cards
}

// Gives a living card to a player.
func (p *Player) GiveCard(card Card) {
	p.Cards = append(p.Cards, CardState{
		Card:  card,
		Alive: true,
	})
}

// Sets the card at `index` to dead.
func (p *Player) KillCard(index int) {
	p.Cards[index].Alive = false
}

// Adjusts the players credits up or down by `amount`.
func (p *Player) AdjustCredits(amount int) {
	p.Credits += amount
}

// Checks if a player can afford an action, erroring if not.
func (p *Player) CanAffordAction(action ActionType) bool {
	// If the action is present in `ActionCost`, it has a cost. Otherwise return true.
	if cost, ok := ActionCost[action]; ok {
		return p.Credits >= cost
	}
	return true
}

// Attempts to deduct the cost of an action from a player's credits, erroring if they cannot afford the action.
func (p *Player) PayForAction(action ActionType) error {
	if cost, ok := ActionCost[action]; ok {
		if p.Credits < cost {
			return errors.New("cannot afford action")
		}
		p.Credits -= cost
	}
	return nil
}

// Tests if the player is allowed to perform an action.
func (p *Player) IsAllowedAction(action ActionType) bool {
	allowed := p.GetAllowedActions()
	return slices.Contains(allowed, action)
}

// Return a list of actions the player is granted by their cards.
func (p *Player) GetAllowedActions() []ActionType {
	granted := DefaultGrants
	for _, card := range p.GetLivingCards() {
		if val, ok := CardGrants[card]; ok {
			granted = append(granted, val)

		}
	}
	return granted
}

// Checks if a player's cards allow them to block
func (p *Player) CanBlock(action ActionType) bool {
	requiredCards, ok := BlockedBy[action]
	// If the action cannot be blocked, return false.
	if !ok {
		return false
	}

	// Check if any of the player's cards allow them to block `action`
	cards := p.GetLivingCards()
	for _, card := range requiredCards {
		if slices.Contains(cards, card) {
			return true
		}
	}
	return false
}
