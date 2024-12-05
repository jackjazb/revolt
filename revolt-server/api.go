package main

import (
	"revolt/game"
)

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
	CreateGameMessage    MessageType = "create_game"
	JoinGameMessage      MessageType = "join_game"
	StartGameMessage     MessageType = "start_game"
	AttemptActionMessage MessageType = "attempt_action"
	CommitTurnMessage    MessageType = "commit_turn"
	EndTurnMessage       MessageType = "end_turn"
)

type CreateGamePayload struct {
	PlayerName string `json:"playerName"`
}

type JoinGamePayload struct {
	GameId     string `json:"gameId"`
	PlayerName string `json:"playerName"`
}

type AttemptActionPayload struct {
	Action game.Action `json:"action"`
}
