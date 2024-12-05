package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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

// Global game instance tracker.
type InstanceManager struct {
	// Maps IDs to game instance pointers (this allows modification)
	Instances map[string]*GameInstance
}

func (im *InstanceManager) RegisterInstance(instance *GameInstance) {
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

	client := NewClient(conn)

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
				delete(currentInstance.Clients, client.Id)
				currentInstance.SendState <- true
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
		case CreateGameMessage:
			var payload CreateGamePayload
			err := UnmarshalPayload(message.Payload, &payload)
			if err != nil {
				client.Log("invalid message payload: %s", message.Payload)
				break
			}

			// Register the instance in the global context.
			instance := NewGameInstance(client.Id)
			im.RegisterInstance(&instance)
			currentInstance = &instance

			// Run handler for client connections and message broadcasts.
			go instance.Run()

			// Register the caller.
			instance.Register <- &client
			instance.SendState <- true
		case JoinGameMessage:
			var payload JoinGamePayload
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
			client.Name = payload.PlayerName
			currentInstance = instance
			instance.Register <- &client
			instance.SendState <- true
		case StartGameMessage:
			if currentInstance == nil {
				client.Log("can't start game (not connected to one)")
				break
			}
			if client.Id != currentInstance.OwnerId {
				client.Log("can't start game (owned by %s)", currentInstance.OwnerId)
				break
			}
			currentInstance.Game.Deal()
			currentInstance.Status = InProgress
			currentInstance.SendState <- true
		case AttemptActionMessage:
			if currentInstance == nil {
				client.Log("can't attempt action - not connected to game")
				break
			}
			var payload AttemptActionPayload
			err = UnmarshalPayload(message.Payload, &payload)
			if err != nil {
				client.Log("error reading message: %s", err)
				break
			}
			err = currentInstance.Game.AttemptAction(payload.Action)
			if err != nil {
				client.Log("couldn't attempt action: %s", err)
				break
			}
			currentInstance.SendState <- true
		case CommitTurnMessage:
			if currentInstance == nil {
				client.Log("can't commit turn - not connected to game")
				break
			}
			err = currentInstance.Game.CommitTurn()
			if err != nil {
				client.Log("couldn't commit action: %s", err)
				break
			}
			currentInstance.SendState <- true
		case EndTurnMessage:
			if currentInstance == nil {
				client.Log("can't end turn - not connected to game")
				break
			}
			err = currentInstance.Game.EndTurn()
			if err != nil {
				client.Log("couldn't end turn: %s", err)
				break
			}
			currentInstance.SendState <- true
		}
	}
}

func RunServer() {
	host := "localhost:8080"
	log.Printf("server up on %s", host)

	// Set up persistent instance tracking struct.
	cm := InstanceManager{
		Instances: make(map[string]*GameInstance),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := MainHandler(&cm, w, r)
		log.Println("error:", err)
	})

	log.Fatal(http.ListenAndServe(host, nil))
}
