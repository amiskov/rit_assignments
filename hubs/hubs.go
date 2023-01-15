package hubs

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type HubsDB struct {
	currentHubID uuid.UUID
	hubs         map[uuid.UUID]*hub
	hubSize      int
}

func NewHubsDB(hubSize int) *HubsDB {
	currentHub := NewHub(hubSize)

	db := HubsDB{
		currentHubID: currentHub.id,
		hubs: map[uuid.UUID]*hub{
			currentHub.id: currentHub,
		},
		hubSize: hubSize,
	}

	log.Printf("HubsDB storage created, current hub is %q.\n", db.hubs[db.currentHubID])
	return &db
}

func (hdb *HubsDB) GetHubById(id string) (*hub, error) {
	hubID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("bad hub id (%w)", err)
	}

	hub, ok := hdb.hubs[hubID]
	if !ok {
		return nil, fmt.Errorf("hub with id %q not found", id)
	}

	return hub, nil
}

func (hdb *HubsDB) GetClientById(id string) (*Client, error) {
	clientID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("bad client id (%w)", err)
	}

	// TODO: Optimise storage to find clients faster.
	for _, hub := range hdb.hubs {
		for _, client := range hub.clients {
			if client.Id == clientID {
				return client, nil
			}
		}
	}

	return nil, fmt.Errorf("client with id %q not found", id)
}

func (hdb HubsDB) ListHubs() []*hub {
	var hubs []*hub
	for _, h := range hdb.hubs {
		hubs = append(hubs, h)
	}
	return hubs
}

func (hdb HubsDB) ListAllClients() []*Client {
	var clients []*Client
	hubs := hdb.ListHubs()
	for _, h := range hubs {
		clients = append(clients, h.clients...)
	}
	return clients
}

func (hdb *HubsDB) Add(c *Client) {
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
