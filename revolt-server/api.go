package main

import (
	"encoding/json"
	"log"
)

// A client state update.
// This should contain everything a client needs to play, but nothing that would allow cheating.
type ClientStateBroadcast struct {
	GameId           string     `json:"gameId"`
	OwnerId          string     `json:"ownerId"`
	ClientId         string     `json:"clientId"`
	Leader           int        `json:"leader"`
	PlayerNumber     int        `json:"number"`
	Status           GameStatus `json:"status"`
	Clients          []Client   `json:"clients"`
	Self             Player     `json:"self"`
	TurnState        TurnState  `json:"turnState"`
	PendingAction    Action     `json:"pendingAction"`
	PendingBlock     Block      `json:"pendingBlock"`
	PendingChallenge Challenge  `json:"pendingChallenge"`
}

func (gi *GameInstance) ToClientStateBroadcast(client *Client) ClientStateBroadcast {
	return ClientStateBroadcast{
		GameId:           gi.GameId,
		OwnerId:          gi.OwnerId,
		ClientId:         client.Id,
		Leader:           gi.Game.Leader,
		PlayerNumber:     client.Number,
		Status:           gi.Status,
		Clients:          []Client{},
		Self:             gi.Game.Players[client.Number],
		TurnState:        gi.Game.TurnState,
		PendingAction:    gi.Game.PendingAction,
		PendingBlock:     gi.Game.PendingBlock,
		PendingChallenge: gi.Game.PendingChallenge,
	}
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

// A message, receivable by the server.
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// Represents a message type.
type MessageType string

// Define possible message types.
const (
	// Game administration messages
	CreateGame MessageType = "create_game"
	JoinGame   MessageType = "join_game"
	StartGame  MessageType = "start_game"
)

// Join game.
type JoinGameMessage struct {
	GameId string `json:"gameId"`
}
