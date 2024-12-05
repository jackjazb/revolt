package game

import (
	"time"

	"golang.org/x/exp/rand"
)

// Defines the different types of card in the game.
type Card string

const (
	Duke       Card = "duke"
	Assassin   Card = "assassin"
	Ambassador Card = "ambassador"
	Captain    Card = "captain"
	Contessa   Card = "contessa"
)

// Defines the state of a single card in a player's hand.
type CardState struct {
	Card  Card `json:"card"`
	Alive bool `json:"alive"`
}

// Defines possible player actions granted by cards.
type ActionType string

const (
	Income     ActionType = "income"
	ForeignAid ActionType = "foreign_aid"
	Revolt     ActionType = "revolt"

	Tax         ActionType = "tax"
	Assassinate ActionType = "assassinate"
	Exchange    ActionType = "exchange"
	Steal       ActionType = "steal"
)

// Defines actions which do not need a card to perform.
var DefaultGrants = []ActionType{Income, ForeignAid, Revolt}

// Defines the actions granted by each type of card.
var CardGrants = map[Card]ActionType{
	Duke:       Tax,
	Assassin:   Assassinate,
	Ambassador: Exchange,
	Captain:    Steal,
}

// Defines which characters can block which actions.
var CardBlocks = map[Card]ActionType{
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

// Check if a card blocks an action.
func (c Card) BlocksAction(action ActionType) bool {
	blocks, ok := CardBlocks[c]
	if ok && blocks == action {
		return true
	}
	return false
}
