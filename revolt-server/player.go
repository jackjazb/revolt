package main

// Defines a single player.
type Player struct {
	Cards   [2]CardState `json:"cards"`
	Credits int          `json:"credits"`
}

// Returns a player's cards.
func (p *Player) GetCards() []Card {
	cards := []Card{}
	for _, card := range p.Cards {
		if card.Alive {
			cards = append(cards, card.Card)
		}
	}
	return cards
}
