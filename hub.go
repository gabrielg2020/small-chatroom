package main

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client
}

// Hub constructor
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run the hub's main loop
// Logic Explanation:
//  1. Continuously listen for messages
//  2. If a there is a value in the register channel, register a new client
//  3. If there is a value in the unregister channel, unregister a client
//  4. If there is a value in the broadcast channel, broadcast a message to all clients
//     4.a If the message is sent successfully, do nothing
//     4.b If the message fails to send, close the client's send channel and delete the client
func (h *Hub) Run() {
	// Continuously listen for messages
	for {
		select {
		case client := <-h.register:
			// Register a new client
			h.clients[client] = true
		case client := <-h.unregister:
			// Unregister a client
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// Broadcast a message to all clients
			for client := range h.clients {
				select {
				case client.send <- message:
					// Message sent successfully
				default:
					// Failed to send message
					// Assume client is disconnected
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
