package main

import (
	"github.com/gorilla/websocket"
)

// client represents a single chatting user.
type client struct {

	// socket is the web socket for this client.
	socket *websocket.Conn

	// send is a channel on which messages are sent.
	send chan []byte

	// room is the room this client is chatting in.
	room *room

	name string
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}

		fullMsg := []byte(c.name + ": " + string(msg))
		c.room.forward <- fullMsg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		//err := c.socket.WriteMessage(websocket.TextMessage, []byte(c.name+": "+string(msg)))
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
