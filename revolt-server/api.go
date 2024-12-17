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
	ChangeNameMessage    MessageType = "change_name"
	StartGameMessage     MessageType = "start_game"
	AttemptActionMessage MessageType = "attempt_action"
	AttemptBlockMessage  MessageType = "attempt_block"
	ChallengeMessage     MessageType = "challenge"
	ResolveDeathMessage  MessageType = "resolve_death"
	CommitTurnMessage    MessageType = "commit_turn"
	EndTurnMessage       MessageType = "end_turn"
)

// Sent to the client on initial connection.
type ConnectionResponse struct {
	Id string `json:"id"`
}

type CreateGamePayload struct {
	PlayerName string `json:"playerName"`
}

type ChangeNamePayload struct {
	PlayerName string `json:"playerName"`
}

type JoinGamePayload struct {
	GameId     string `json:"gameId"`
	PlayerName string `json:"playerName"`
}
type RejoinGamePayload struct {
	GameId   string `json:"gameId"`
	ClientId string `json:"clientId"`
}

type AttemptActionPayload struct {
	Action game.Action `json:"action"`
}

type AttemptBlockPayload struct {
	Block game.Block `json:"block"`
}

type ChallengePayload struct {
	Challenge game.Challenge `json:"challenge"`
}

type ResolveDeathPayload struct {
	Card int `json:"card"`
}
