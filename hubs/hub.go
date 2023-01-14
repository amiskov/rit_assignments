package hubs

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type hubID uuid.UUID

type hub struct {
	id      hubID
	size    int
	clients []*client
}

func NewHub(size int) *hub {
	return &hub{
		id:      hubID(uuid.New()),
		size:    size,
		clients: []*client{},
	}
}

// TODO: We also have to handle WS disconnect for hubs
func (h *hub) Append(client *client) error {
	if len(h.clients) >= h.size {
		fmt.Println(h, h.size)
		return ErrLimitExceeded
	}
	h.clients = append(h.clients, client)
	log.Printf("client %s was added to hub %s\n", client, h)
	return nil
}

// Sends a message to all clients of a hub
func (h *hub) broadcast() {
	for idx, c := range h.clients {
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			panic("client failed")
		}
		msg := fmt.Sprintf("sent from hub %d to client %d", -1, idx)
		w.Write([]byte(msg))
		w.Close()
	}
}

func (h *hub) String() string {
	return uuid.UUID(h.id).String()
}
