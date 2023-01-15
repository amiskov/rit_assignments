package hubs

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

type hubID uuid.UUID

type hub struct {
	id      hubID
	size    int // TODO: move size checking into HubsDB
	clients []*Client
}

func NewHub(size int) *hub {
	return &hub{
		id:      hubID(uuid.New()),
		size:    size,
		clients: []*Client{},
	}
}

// TODO: We also have to handle WS disconnect for hubs
func (h *hub) Append(client *Client) error {
	// TODO: move size checking into HubsDB
	if len(h.clients) >= h.size {
		fmt.Println(h, h.size)
		return ErrLimitExceeded
	}
	h.clients = append(h.clients, client)
	log.Printf("client %s was added to hub %s\n", client, h)
	return nil
}

func (h *hub) ListClients() []*Client {
	return h.clients
}

// Sends a message to all clients of a hub
func (h *hub) Broadcast(msg string) error {
	var wg sync.WaitGroup
	for _, c := range h.clients {
		wg.Add(1)
		err := c.SendMessage(msg)
		if err != nil {
			return fmt.Errorf("failed broadcasting message to the Hub %q clients (%w)", h, err)
		}
		wg.Done()
	}
	wg.Wait()
	log.Printf("Broadcasted message %q to %d clients of the hub %q\n", msg, len(h.clients), h)
	return nil
}

func (h *hub) String() string {
	return uuid.UUID(h.id).String()
}
