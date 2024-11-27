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
	Id        string
	Status    GameStatus
	GameState Game
	Clients   map[string]*Client
	// Channel to register new clients with the game instance.
	Register chan *Client
	// Channel to send messages to all the instance's clients.
	Broadcast chan *Message
}

func NewGameInstance() GameInstance {
	return GameInstance{
		Id:        Id(),
		Status:    Lobby,
		Clients:   make(map[string]*Client),
		GameState: Game{},
		Register:  make(chan *Client),
		Broadcast: make(chan *Message),
	}
}

// Receives client registrations and broadcast messages on their respective channels and handles them.
func (c *GameInstance) Run() {
	log.Printf("running new game instance %s", c.Id)
	for {
		select {
		case client := <-c.Register:
			log.Printf("registering client %s with game %s", client.Id, c.Id)
			c.Clients[client.Id] = client
		case message := <-c.Broadcast:
			log.Printf("broadcasting message %+v to game %s", message, c.Id)
			bytes, err := message.Serialise()
			if err != nil {
				break
			}
			for _, client := range c.Clients {
				client.Send <- bytes
			}
		}
	}
}

// Tracks game instances.
type InstanceManager struct {
	Instances map[string]GameInstance
}

func (im *InstanceManager) RegisterInstance(instance GameInstance) {
	im.Instances[instance.Id] = instance
}

// Represents a client
type Client struct {
	Id         string
	Connection *websocket.Conn
	Send       chan []byte
}

// WebSocket handler.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Sends messages received on the `Send` channel to connected the client.
func (c *Client) HandleMessages() {
	defer c.Connection.Close()
	for message := range c.Send {
		if err := c.Connection.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err)
		}
	}
}

// Primary websocket connection handler.
func MainHandler(im *InstanceManager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := Client{
		Id:         Id(),
		Connection: conn,
		Send:       make(chan []byte),
	}
	go client.HandleMessages()

	// Register the client with
	// cm.Register <- &client

	log.Printf("new client connected: %s", client.Id)

	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Parse JSON contents of message.
		var message Message
		err = json.Unmarshal(bytes, &message)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("received message %+v", message)

		switch message.Type {
		case CreateGame:
			// Register the instance in the global context.
			instance := NewGameInstance()
			im.RegisterInstance(instance)

			// Run handler for client connections and message broadcasts.
			go instance.Run()

			// Register the caller.
			instance.Register <- &client
			instance.Broadcast <- &Message{
				Type: CreateGame,
				Payload: CreateGameResponse{
					GameId:   instance.Id,
					ClientId: client.Id,
				},
			}
		case JoinGame:
			var payload JoinGameMessage
			err := UnmarshalPayload(message.Payload, &payload)
			if err != nil {
				log.Printf("invalid message payload: %s", message.Payload)
				break
			}

			instance, ok := im.Instances[payload.GameId]
			if !ok {
				log.Printf("instance with id %s not found", payload.GameId)
				break
			}
			instance.Register <- &client
			// TODO store player index on client
			instance.Broadcast <- &Message{
				Type: JoinGame,
				Payload: JoinGameResponse{
					ClientId: client.Id,
				},
			}
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

	// go cm.Manage()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		MainHandler(&cm, w, r)
	})
	log.Fatal(http.ListenAndServe(host, nil))
}
