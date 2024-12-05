package main

import (
	"encoding/json"
	"log"
	"revolt/game"
)

// The status of a given instance.
type GameStatus string

// Defines possible instance statuses.
const (
	Lobby              GameStatus = "lobby"
	InProgress         GameStatus = "in_progress"
	CompleteGameStatus GameStatus = "complete"
)

// A single instance of a game.
type GameInstance struct {
	GameId  string
	OwnerId string
	Status  GameStatus
	Game    game.Game
	Clients map[string]*Client

	// TODO unregister clients as well.
	Register  chan *Client // Channel to register new clients with the game instance.
	SendState chan bool    // Channel to trigger a state broadcast
}

// Creates a new game instance in the `Lobby` status.
func NewGameInstance(ownerId string) GameInstance {
	return GameInstance{
		GameId:    game.Id(),
		OwnerId:   ownerId,
		Status:    Lobby,
		Clients:   make(map[string]*Client),
		Game:      game.NewGame(),
		Register:  make(chan *Client),
		SendState: make(chan bool),
	}
}

// Receives client registrations on the Register channel and handles state broadcast requests.
func (gi *GameInstance) Run() {
	log.Printf("running new game instance %s", gi.GameId)
	for {
		select {
		// Registers a client with the current game instance.
		case client := <-gi.Register:
			client.Log("registering client %s with game %s...", client.Name, gi.GameId)

			// Add the player to the current game instance.
			err := gi.Game.AddPlayer(client.Id, client.Name)
			if err != nil {
				client.Log("error registering client with game instance")
				continue
			}

			gi.Clients[client.Id] = client

		// Triggers a broadcast of the current instance state to all connected clients.
		case <-gi.SendState:
			log.Printf("broadcasting state to game instance %s", gi.GameId)

			for _, client := range gi.Clients {
				update := gi.ToClientStateBroadcast(client)
				bytes, err := update.Serialise()
				if err != nil {
					break
				}
				client.Send <- bytes
			}
		}
	}
}

// A client state update.
// This should contain everything a client needs to play, but nothing that would allow cheating.
type ClientStateBroadcast struct {
	// Session and client info.
	GameId  string     `json:"gameId"`
	OwnerId string     `json:"ownerId"`
	Self    Peer       `json:"self"`
	Peers   []Peer     `json:"peers"`
	Status  GameStatus `json:"status"`

	// Game info.
	Leader           int            `json:"leader"`
	PlayerNumber     int            `json:"number"`
	TurnState        game.TurnState `json:"turnState"`
	PendingAction    game.Action    `json:"pendingAction"`
	PendingBlock     game.Block     `json:"pendingBlock"`
	PendingChallenge game.Challenge `json:"pendingChallenge"`
}

type Peer struct {
	Name    string           `json:"name"`
	Id      string           `json:"id"`
	Cards   []game.CardState `json:"cards"`
	Credits int              `json:"credits"`
	Leading bool             `json:"leading"`
}

// Converts a state update message to a JSON byte array.
func (s *ClientStateBroadcast) Serialise() ([]byte, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return bytes, nil
}

func (gi *GameInstance) ToClientStateBroadcast(client *Client) ClientStateBroadcast {

	// Collect relevant information
	peers := []Peer{}
	self := Peer{}
	for i, player := range gi.Game.Players {
		dead := player.GetDeadCards()
		peer := Peer{
			Id:      player.Id,
			Name:    player.Name,
			Cards:   dead,
			Credits: player.Credits,
			Leading: i == gi.Game.Leader,
		}
		if player.Id == client.Id {
			self = peer
			continue
		}
		// We only want to send other player's dead cards.
		peers = append(peers, peer)
	}

	return ClientStateBroadcast{
		GameId:  gi.GameId,
		OwnerId: gi.OwnerId,
		Self:    self,
		Peers:   peers,
		Status:  gi.Status,

		Leader:           gi.Game.Leader,
		TurnState:        gi.Game.TurnState,
		PendingAction:    gi.Game.PendingAction,
		PendingBlock:     gi.Game.PendingBlock,
		PendingChallenge: gi.Game.PendingChallenge,
	}
}
