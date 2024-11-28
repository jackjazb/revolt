package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Represents the state of a connected client.
type Client struct {
	Id         string          `json:"id"`
	Name       string          `json:"name"`
	Number     int             `json:"number"`
	Connection *websocket.Conn `json:"-"`
	Send       chan []byte     `json:"-"`
}

// Sends messages received on the `Send` channel to connected the client.
func (c *Client) HandleMessages() {
	defer c.Connection.Close()
	for message := range c.Send {
		if err := c.Connection.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err)
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
