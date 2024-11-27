package main

import (
	"encoding/json"
	"log"
)

// A message, receivable by the server. `T` defines the payload type.
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// Converts a message struct to a JSON byte array.
func (m *Message) Serialise() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return bytes, nil
}

// Represents a message type.
type MessageType string

// Define possible message types.
const (
	CreateGame MessageType = "create_game"
	JoinGame   MessageType = "join_game"
)

// Create game.
type CreateGameResponse struct {
	GameId   string `json:"gameId"`
	ClientId string `json:"clientId"`
}

// Join game.
type JoinGameMessage struct {
	GameId string `json:"gameId"`
}

type JoinGameResponse struct {
	ClientId string `json:"clientId"`
}
