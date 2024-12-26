package main

import (
	"fmt"
	"log"
	"revolt/game"

	"github.com/gorilla/websocket"
)

// Represents the state of a connected client.
type Client struct {
	Id         string
	Name       string
	Connection *websocket.Conn
	Send       chan []byte
}

func NewClient(conn *websocket.Conn, name string) Client {
	return Client{
		Id:         game.Id(),
		Name:       name,
		Connection: conn,
		Send:       make(chan []byte),
	}
}

func (c *Client) HandleMessages() {
	defer c.Connection.Close()
	for message := range c.Send {
		if err := c.Connection.WriteMessage(websocket.TextMessage, message); err != nil {
			c.Log("error writing message: %s", err)
			return
		}
	}
	c.Log("client handler stopped")
}

// Utility function for logging events that happen in the context of a client.
func (c *Client) Log(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.Printf("%s (client: %s)", message, c.Id)
}
