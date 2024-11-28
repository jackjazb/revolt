package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func Id() string {
	// TODO this should be a full length UUID, but 8 chars is easier to read in development.
	return uuid.NewString()[:8]
}

// Unmarshals a payload interface, allowing type assignment.
func UnmarshalPayload(payload interface{}, v any) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		return err
	}
	return nil
}

// The status of a given instance.
type GameStatus int

// Defines possible instance statuses.
const (
	Lobby GameStatus = iota
	InProgress
	Complete
)

// A single instance of a game.
type GameInstance struct {
	GameId  string             `json:"gameId"`
	OwnerId string             `json:"ownerId"`
	Status  GameStatus         `json:"status"`
	Game    Game               `json:"game"`
	Clients map[string]*Client `json:"clients"`

	// TODO unregister clients as well.
	Register  chan *Client `json:"-"` // Channel to register new clients with the game instance.
	SendState chan bool    `json:"-"` // Channel to trigger a state broadcast
}

func NewGameInstance(ownerId string) GameInstance {
	return GameInstance{
		GameId:    Id(),
		OwnerId:   ownerId,
		Status:    Lobby,
		Clients:   make(map[string]*Client),
		Game:      NewGame(),
		Register:  make(chan *Client),
		SendState: make(chan bool),
	}
}

// Receives client registrations and broadcast messages on their respective channels and handles them.
func (gi *GameInstance) Run() {
	log.Printf("running new game instance %s", gi.GameId)
	for {
		select {

		// Registers a client with the current game instance.
		case client := <-gi.Register:
			client.Log("registering client with game %s...", gi.GameId)

			// Add the player to the current game instance.
			number, err := gi.Game.AddPlayer()
			if err != nil {
				client.Log("error registering client with game instance")
				continue
			}

			client.Number = number
			gi.Clients[client.Id] = client

		// Triggers a broadcast of the current instance state to all connected clients.
		case <-gi.SendState:
			log.Printf("broadcasting state to game instance %s", gi.GameId)

			for _, client := range gi.Clients {
				update := &StateUpdate{
					ClientId: client.Id,
					State:    *gi,
				}
				bytes, err := update.Serialise()
				if err != nil {
					break
				}
				client.Send <- bytes
			}
		}
	}
}

// Global game instance tracker.
type InstanceManager struct {
	Instances map[string]GameInstance
}

func (im *InstanceManager) RegisterInstance(instance GameInstance) {
	im.Instances[instance.GameId] = instance
}

// WebSocket handler.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Primary websocket connection handler.
func MainHandler(im *InstanceManager, w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	client := Client{
		Id:         Id(),
		Connection: conn,
		Send:       make(chan []byte),
	}

	// Listen for messages and write them when received.
	go client.HandleMessages()

	log.Printf("new client connected: %s", client.Id)

	// Holds a pointer to the client's current game instance.
	var currentInstance *GameInstance

	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			// If the client closed the connection, make sure to stop their handler.
			ce, ok := err.(*websocket.CloseError)
			if !ok {
				return err
			}
			switch ce.Code {
			case websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived:
				client.Log("connection closed by client")
				close(client.Send)
				return err
			}
		}

		// Parse the received message.
		var message Message
		err = json.Unmarshal(bytes, &message)
		if err != nil {
			return err
		}

		log.Printf("received message %+v", message)

		switch message.Type {
		case CreateGame:
			// Register the instance in the global context.
			instance := NewGameInstance(client.Id)
			im.RegisterInstance(instance)
			currentInstance = &instance
			// Run handler for client connections and message broadcasts.
			go instance.Run()

			// Register the caller.
			instance.Register <- &client
			instance.SendState <- true
		case JoinGame:
			var payload JoinGameMessage
			err := UnmarshalPayload(message.Payload, &payload)
			if err != nil {
				client.Log("invalid message payload: %s", message.Payload)
				break
			}

			instance, ok := im.Instances[payload.GameId]
			if !ok {
				client.Log("instance with id %s not found", payload.GameId)
				break
			}
			currentInstance = &instance
			instance.Register <- &client
			instance.SendState <- true
		case StartGame:
			if currentInstance == nil {
				client.Log("client not currently connected to a game")
				break
			}
			if client.Id != currentInstance.OwnerId {
				client.Log("can't start game - owned by %s", currentInstance.OwnerId)
				break
			}
			currentInstance.Game.Deal()
			currentInstance.Status = InProgress
			currentInstance.SendState <- true
		}
	}
}

func RunServer() {
	host := "localhost:8080"
	log.Printf("server up on %s", host)

	// Set up persistent instance tracking struct.
	cm := InstanceManager{
		Instances: make(map[string]GameInstance),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := MainHandler(&cm, w, r)
		log.Println("error:", err)
	})

	log.Fatal(http.ListenAndServe(host, nil))
}
