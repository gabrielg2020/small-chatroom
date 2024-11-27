package main

import (
	"bytes"
  "log"
	"time"

	"github.com/gorilla/websocket"
)

// Time related constants
const (
	waitTime = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Predefined byte slices for message processing
var (
	newline = []byte{'\n'}
	space = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The hub
	hub *Hub

	// The websocket connection
	conn *websocket.Conn

	// Channel to send messages to the hub
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump() {
	defer func() {
		// On exit, unregister the client and close the connection
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Set the max message size and initial read deadline
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	// Define the pong handler to reset the read deadline
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Continuously read messages from the connection
	for {
		// Read a message from the connection
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// Log the error
				log.Printf("error: %v", err)
			}
			break
		}

		// Clean up the message
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// Send the message to the hub
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	// Create a ticker to send ping messages
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		// On exit, stop the ticker and close the connection
		ticker.Stop()
		c.conn.Close()
	}()

	// Continuously listen for messages
	for {
		select {
		case message, ok := <-c.send:
			// Write a message to the connection
			c.conn.SetWriteDeadline(time.Now().Add(waitTime))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Write the message to the WebSocket connection
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
	 		}

			// Close the writer
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// Send a ping message
			c.conn.SetWriteDeadline(time.Now().Add(waitTime))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}