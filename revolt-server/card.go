package main

import "golang.org/x/exp/rand"

// Defines types of card,
type Card int

const (
	Duke Card = iota
	Assassin
	Ambassador
	Captain
	Contessa
)

// Defines the state of a single player card.
type CardState struct {
	Card  Card
	Alive bool
}

// Defines possible player actions.
type ActionType int

const (
	Income ActionType = iota
	ForeignAid
	Tax
	Assassinate
	Revolt
	Exchange
	Steal
)

// Defines which cards grant which actions.
var RequiredCard = map[ActionType]Card{
	Tax:         Duke,
	Assassinate: Assassin,
	Exchange:    Ambassador,
	Steal:       Captain,
}

// Defines which characters can block which actions.
var Blocks = map[Card]ActionType{
	Duke:       ForeignAid,
	Contessa:   Assassinate,
	Captain:    Steal,
	Ambassador: Steal,
}

// Defines which actions are blocked by which characters (inverse of Blocks).
var BlockedBy = map[ActionType][]Card{
	ForeignAid:  {Duke},
	Assassinate: {Contessa},
	Steal:       {Captain, Ambassador},
}

// Those cost of different action types.
var ActionCost = map[ActionType]int{
	Assassinate: 3,
	Revolt:      7,
}

// The starting deck.
var Deck = []Card{
	Duke, Assassin, Ambassador, Captain, Contessa,
	Duke, Assassin, Ambassador, Captain, Contessa,
	Duke, Assassin, Ambassador, Captain, Contessa,
}

// Formats a card.
func (i Card) String() string {
	return []string{
		"Duke",
		"Assassin",
		"Ambassador",
		"Captain",
		"Contessa",
	}[i]
}

// Shuffles a list of cards.
func ShuffleCards(cards []Card) []Card {
	shuffled := cards
	for range 16 {
		deck := shuffled
		shuffled = []Card{}
		for len(deck) > 1 {
			lastIndex := len(deck) - 1
			randomIndex := rand.Intn(lastIndex)
			shuffled = append(shuffled, deck[randomIndex])

			// Replace randomIndex with the last element and remove the last element.
			deck[randomIndex] = deck[lastIndex]
			deck = deck[:lastIndex]
		}
		shuffled = append(shuffled, deck[0])
	}
	return shuffled
}
