package main

import (
	"encoding/json"
	"log"
)

// The response from the server for client actions, for broadcasting the current state.
type StateUpdate struct {
	ClientId string       `json:"clientId"`
	State    GameInstance `json:"state"`
}

// Converts a state update message to a JSON byte array.
func (s *StateUpdate) Serialise() ([]byte, error) {
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
