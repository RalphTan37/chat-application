package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket  *websocket.Conn //socket is the web socket for the client
	recieve chan []byte     //recieve channel - to recieve messages from other clients
	room    *room           //room the client is chatting in
}

// read method
func (c *client) read() {
	defer c.socket.Close() //whenever read function returns, close the connection
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg //forward message to room when recieved from socket
	}
}

// write method
func (c *client) write() {
	defer c.socket.Close()       // whenever write functionr returns, close the connection
	for msg := range c.recieve { //read message from client's recieve channel
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
