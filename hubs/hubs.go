package hubs

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type HubsDB struct {
	currentHubID hubID
	hubs         map[hubID]*hub
	hubSize      int
}

func NewHubsDB(hubSize int) *HubsDB {
	currentHub := NewHub(hubSize)

	db := HubsDB{
		currentHubID: currentHub.id,
		hubs: map[hubID]*hub{
			currentHub.id: currentHub,
		},
		hubSize: hubSize,
	}

	log.Printf("HubsDB storage created, current hub is %q.\n", db.hubs[db.currentHubID])
	return &db
}

func (hdb *HubsDB) GetHubById(id string) (*hub, error) {
	hUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("bad hub id (%w)", err)
	}

	hub, ok := hdb.hubs[hubID(hUUID)]
	if !ok {
		return nil, fmt.Errorf("hub with id %q not found", id)
	}

	return hub, nil
}

func (hdb *HubsDB) GetClientById(id string) (*client, error) {
	cUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("bad client id (%w)", err)
	}

	cID := clientID(cUUID)

	// TODO: Optimise storage to find clients faster.
	for _, hub := range hdb.hubs {
		for _, client := range hub.clients {
			if client.id == cID {
				return client, nil
			}
		}
	}

	return nil, fmt.Errorf("client with id %q not found", id)
}

func (hdb *HubsDB) Add(ws *websocket.Conn) {
	client := client{
		id:   clientID(uuid.New()),
		conn: ws,
	}

	// Try to save to current hub
	currentHub := hdb.hubs[hdb.currentHubID]
	err := currentHub.Append(&client)
	if errors.Is(err, ErrLimitExceeded) {
		newHub := NewHub(hdb.hubSize)
		newHub.Append(&client)
		hdb.hubs[newHub.id] = newHub
		hdb.currentHubID = newHub.id
		return
	}

	if err != nil {
		// TODO: improve error
		log.Fatalln("failed to add new client to hub:", err)
	}
}
