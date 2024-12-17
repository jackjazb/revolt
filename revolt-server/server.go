package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func remove(array []string, value string) (ret []string) {
	for _, s := range array {
		if s != value {
			ret = append(ret, s)
		}
	}
	return
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

	bytes, err := json.Marshal(ConnectionResponse{Id: client.Id})
	if err != nil {
		client.Log("connection failed immediately, %v", err)
		return err
	}
	// Send the client's ID on initial connection.
	client.Send <- bytes

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

				return err
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
		// TODO we could do this automatically on first load.
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

			// Set the owner to the first client to connect.
			if len(instance.Clients) == 0 {
				instance.OwnerId = client.Id
			}

			client.Name = payload.PlayerName
			currentInstance = instance
			instance.Register <- &client
			instance.SendState <- true

		case ChangeNameMessage:
			if currentInstance == nil {
				client.Log("can't change player name - not connected to game")
				break
			}
			var payload ChangeNamePayload
			err := UnmarshalPayload(message.Payload, &payload)
			if err != nil {
				client.Log("invalid message payload: %s", message.Payload)
				break
			}

			client.Name = payload.PlayerName
			currentInstance.SendState <- true

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

func RunServer() {
	host := "localhost:8080"
	log.Printf("server up on %s", host)

	// Set up persistent instance tracking struct.
	im := InstanceManager{
		Instances: make(map[string]*GameInstance),
	}

	http.HandleFunc("POST /create", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Register the instance in the global context.
		// TODO set owner on first client connection
		instance := NewGameInstance("")
		im.RegisterInstance(&instance)

		// Run handler for client connections and message broadcasts.
		go instance.Run()

		bytes, err := json.Marshal(ConnectionResponse{Id: instance.GameId})
		if err != nil {
			log.Println("failed to send id")
			return
		}
		w.Write(bytes)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := MainHandler(&im, w, r)
		log.Println("error in main handler:", err)
	})

	log.Fatal(http.ListenAndServe(host, nil))
}
