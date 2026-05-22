package ws

import (
	"encoding/json"
)

type Event struct {
	RunID   string
	Topic   string
	Payload interface{}
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Event
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan Event),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
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
		case event := <-h.broadcast:
			payload, err := json.Marshal(event.Payload)
			if err != nil {
				continue
			}

			for client := range h.clients {
				if client.runID == event.RunID && client.topic == event.Topic {
					select {
					case client.send <- payload:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

func (h *Hub) Broadcast(runID, topic string, payload interface{}) {
	h.broadcast <- Event{
		RunID:   runID,
		Topic:   topic,
		Payload: payload,
	}
}
