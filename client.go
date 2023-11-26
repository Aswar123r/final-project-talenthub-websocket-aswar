package main

import (
	"github.com/gorilla/websocket"
)

type ChatClient struct {
	socket  *websocket.Conn
	receive chan []byte
	room    *ChatRoom
}

func (c *ChatClient) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *ChatClient) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
