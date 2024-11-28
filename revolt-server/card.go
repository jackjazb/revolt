package main

import (
	"time"

	"golang.org/x/exp/rand"
)

// Defines types of card,
type Card string

const (
	Duke       Card = "duke"
	Assassin   Card = "assassin"
	Ambassador Card = "ambassador"
	Captain    Card = "captain"
	Contessa   Card = "contessa"
)

// Defines the state of a single player card.
type CardState struct {
	Card  Card `json:"card"`
	Alive bool `json:"alive"`
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

// Shuffles a list of cards.
func ShuffleCards(cards []Card) []Card {
	rand.Seed(uint64(time.Now().UnixNano()))
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
