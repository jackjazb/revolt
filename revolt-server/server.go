package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

const NameKey = "name"

func remove(array []string, value string) (ret []string) {
	for _, s := range array {
		if s != value {
			ret = append(ret, s)
		}
	}
	return
}

// Writes an error message to a websocket connection and closes it.
func errorAndClose(conn *websocket.Conn, error string) {
	conn.WriteMessage(websocket.CloseMessage, []byte(fmt.Sprintf(`{"error":"%s"}`, error)))
	conn.Close()
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

// Struct for tracking instances
type InstanceManager struct {
	// Maps IDs to game instance pointers (this allows modification)
	Instances map[string]*GameInstance
}

// Global instance store.
var im InstanceManager

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
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new websocket connection.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Check the path for an instance ID.
	path := strings.Split(r.URL.Path, "/")[1:]
	if len(path) != 1 {
		errorAndClose(conn, "missing ID in URL")
		return
	}

	id := path[0]
	instance, ok := im.Instances[id]
	if !ok {
		errorAndClose(conn, "instance not found")
		return
	}

	// Extract a name from the URL if present
	name := r.URL.Query().Get(NameKey)
	client := NewClient(conn, name)

	go client.HandleMessages()
	client.Log("new client connected with name %s", client.Name)

	// Holds a pointer to the client's current game instance.
	var currentInstance *GameInstance

	// Set the owner to the first client to connect.
	if len(instance.Clients) == 0 {
		instance.OwnerId = client.Id
	}

	currentInstance = instance
	instance.Register <- &client
	instance.SendState <- true

	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			// If the client closed the connection, make sure to stop their handler.
			ce, ok := err.(*websocket.CloseError)
			if !ok {
				log.Println(err)
				return
			}
			switch ce.Code {
			case
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived:

				client.Log("connection closed by client")

				if currentInstance != nil {
					delete(currentInstance.Clients, client.Id)
					delete(currentInstance.Game.Players, client.Id)
					currentInstance.Game.Order = remove(currentInstance.Game.Order, client.Id)
					currentInstance.SendState <- true
				}
				close(client.Send)

				log.Println(err)
				return
			}
		}

		// Parse the received message.
		var message Message
		err = json.Unmarshal(bytes, &message)
		if err != nil {
			client.Log("invalid message payload: %s", message.Payload)
			continue
		}

		log.Printf("received message %+v", message)

		switch message.Type {
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

			// TODO remove, only for debug purposes.
			for _, p := range currentInstance.Game.Players {
				p.Credits += 5
			}

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

		case AttemptBlockMessage:
			if currentInstance == nil {
				client.Log("can't attempt block - not connected to game")
				break
			}
			var payload AttemptBlockPayload
			err = UnmarshalPayload(message.Payload, &payload)
			if err != nil {
				client.Log("error reading message: %s", err)
				break
			}
			// Set initiator - even if provided, we don't want to allow impersonating other players.
			payload.Block.Initiator = client.Id
			err = currentInstance.Game.AttemptBlock(payload.Block)
			if err != nil {
				client.Log("couldn't attempt block: %s", err)
				break
			}
			currentInstance.SendState <- true

		case ChallengeMessage:
			if currentInstance == nil {
				client.Log("can't attempt block - not connected to game")
				break
			}
			var payload ChallengePayload
			err = UnmarshalPayload(message.Payload, &payload)
			if err != nil {
				client.Log("error reading message: %s", err)
				break
			}
			// Set initiator - even if provided, we don't want to allow impersonating other players.
			payload.Challenge.Initiator = client.Id
			err = currentInstance.Game.Challenge(payload.Challenge)
			if err != nil {
				client.Log("couldn't attempt challenge: %s", err)
				break
			}
			currentInstance.SendState <- true

		case ResolveDeathMessage:
			if currentInstance == nil {
				client.Log("can't resolve death - not connected to game")
				break
			}
			var payload ResolveDeathPayload
			err = UnmarshalPayload(message.Payload, &payload)
			if err != nil {
				client.Log("error reading message: %s", err)
				break
			}
			err = currentInstance.Game.ResolveDeath(payload.Card)
			if err != nil {
				client.Log("couldn't resolve death action: %s", err)
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

func createGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not permitted", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Register the instance in the global context.
	instance := NewGameInstance("")
	im.RegisterInstance(&instance)

	// Run handler for client connections and message broadcasts.
	go instance.Run()

	bytes, err := json.Marshal(ConnectionResponse{Id: instance.GameId})
	if err != nil {
		log.Println("failed to send id of new game")
		return
	}
	w.Write(bytes)
}

func initInstanceManager() {
	im = InstanceManager{
		Instances: make(map[string]*GameInstance),
	}
}

func RunServer() error {
	host := "localhost:8080"
	log.Printf("server up on %s", host)

	initInstanceManager()

	mux := http.NewServeMux()
	mux.Handle("/create", http.HandlerFunc(createGameHandler))
	mux.Handle("/{id}", http.HandlerFunc(websocketHandler))

	err := http.ListenAndServe(host, mux)
	if err != nil {
		return err
	}
	return nil
}
