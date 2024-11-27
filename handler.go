package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all connection origins... Note: This is not secure
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWs handles websocket requests from the client.
// Logic Explanation:
//  1. Upgrade the HTTP connection to a websocket connection
//  2. Create a new client
//  3. Register the client with the hub
//  4. Start the clients write and read pumps
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create a new client
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}

	// Register the client with the hub
	client.hub.register <- client

	// Start the clients write and read pumps
	go client.WritePump()
	go client.ReadPump()
}
