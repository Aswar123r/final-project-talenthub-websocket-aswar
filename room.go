package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type ChatRoom struct {
	clients map[*ChatClient]bool
	join    chan *ChatClient
	leave   chan *ChatClient
	forward chan []byte
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		forward: make(chan []byte),
		join:    make(chan *ChatClient),
		leave:   make(chan *ChatClient),
		clients: make(map[*ChatClient]bool),
	}
}

func (r *ChatRoom) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.forward:
			for client := range r.clients {
				client.receive <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *ChatRoom) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &ChatClient{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
