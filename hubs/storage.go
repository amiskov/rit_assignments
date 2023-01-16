package hubs

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

type HubsStorage struct {
	mx sync.RWMutex

	currentHubID uuid.UUID
	hubs         map[uuid.UUID]*Hub
	hubSize      int
}

func NewHubsStorage(hubSize int) *HubsStorage {
	currentHub := NewHub(hubSize)

	db := HubsStorage{
		currentHubID: currentHub.id,
		hubs: map[uuid.UUID]*Hub{
			currentHub.id: currentHub,
		},
		hubSize: hubSize,
	}

	log.Printf("HubsDB storage created, current hub is %q.\n", db.hubs[db.currentHubID])
	return &db
}

func (hdb *HubsStorage) GetHubById(id string) (*Hub, error) {
	hubID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("bad hub id (%w)", err)
	}

	hdb.mx.RLock()
	hub, ok := hdb.hubs[hubID]
	if !ok {
		return nil, fmt.Errorf("hub with id %q not found", id)
	}
	hdb.mx.RUnlock()

	return hub, nil
}

func (hdb *HubsStorage) GetClientById(id string) (*Client, error) {
	clientID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("bad client id (%w)", err)
	}

	// TODO: Optimise storage to find clients faster.
	hdb.mx.RLock()
	for _, hub := range hdb.hubs {
		for _, client := range hub.clients {
			if client.Id == clientID {
				return client, nil
			}
		}
	}
	hdb.mx.RUnlock()

	return nil, fmt.Errorf("client with id %q not found", id)
}

func (hdb *HubsStorage) ListHubs() []*Hub {
	var hubs []*Hub
	hdb.mx.Lock()
	defer hdb.mx.Unlock()
	for _, h := range hdb.hubs {
		hubs = append(hubs, h)
	}
	return hubs
}

func (hdb *HubsStorage) ListAllClients() []*Client {
	var clients []*Client
	hubs := hdb.ListHubs()

	hdb.mx.RLock()
	defer hdb.mx.RUnlock()
	for _, h := range hubs {
		clients = append(clients, h.clients...)
	}
	return clients
}

func (hdb *HubsStorage) Add(c *Client) {
	hdb.mx.Lock()
	defer hdb.mx.Unlock()
	currentHub := hdb.hubs[hdb.currentHubID]

	if len(currentHub.clients) < hdb.hubSize {
		currentHub.Add(c)
		return
	}

	newHub := NewHub(hdb.hubSize)
	newHub.Add(c)
	hdb.hubs[newHub.id] = newHub
	hdb.currentHubID = newHub.id
}
