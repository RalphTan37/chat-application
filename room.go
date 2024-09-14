package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	clients map[*client]struct{} //all current clients in the room
	join    chan *client         //channel for clients to join the room
	leave   chan *client         //channel for clients to leave the room
	forward chan []byte          //channel that holds incoming messages that'll be forwarded to the other clients
}

// constructor for the room
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]struct{}),
	}
}

// method to run the room
func (r *room) run() {
	//listens for values on the channels w/in the room
	for {
		select {
		case client := <-r.join:
			r.clients[client] = struct{}{}
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.recieve)
		case msg := <-r.forward:
			for client := range r.clients {
				client.recieve <- msg
			}
		}
	}
}

// amount of data the buffers can hold before reading & writing values/calls to the network stack
const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServerHTTP(w http.ResponseWriter, req *http.Request) {
	//upgrades the connection
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServerHTTP", err)
		return
	}

	//create a new instance of a client
	client := &client{
		socket:  socket,
		recieve: make(chan []byte, messageBufferSize),
		room:    r,
	}
	r.join <- client
	defer func() { r.leave <- client }() //if room closes, remove client
	go client.write()
	client.read()
}
