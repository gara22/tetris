package socket

import "fmt"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcastChannel chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	IncomingChannel chan []byte
}

func NewHub() *Hub {
	return &Hub{
		broadcastChannel: make(chan []byte),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		clients:          make(map[*Client]bool),
		IncomingChannel:  make(chan []byte),
	}
}

func (h Hub) PublishMessage(data []byte) {
	fmt.Println("Publishing message")
	h.broadcastChannel <- data
}

func (h Hub) ReadMessage() chan []byte {
	return h.IncomingChannel
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcastChannel:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
