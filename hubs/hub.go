package hubs

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	id      uuid.UUID
	clients []*Client
}

func NewHub(size int) *Hub {
	return &Hub{
		id:      uuid.New(),
		clients: []*Client{},
	}
}

func (h *Hub) Add(client *Client) {
	h.clients = append(h.clients, client)
	log.Printf("client %s was added to hub %s\n", client, h)
}

func (h *Hub) ListClients() []*Client {
	return h.clients
}

// Sends a message to all clients of a hub
func (h *Hub) Broadcast(msg string) {
	var wg sync.WaitGroup
	for _, c := range h.clients {
		wg.Add(1)
		go func(c *Client) {
			defer wg.Done()
			err := c.SendMessage(msg)
			if err != nil {
				log.Printf("failed sending message to the client %s (%v)", c, err)
			}
		}(c)
	}
	wg.Wait()
	log.Printf("Broadcasted message %q to %d clients of the hub %q\n", msg, len(h.clients), h)
}

func (h *Hub) String() string {
	return uuid.UUID(h.id).String()
}
