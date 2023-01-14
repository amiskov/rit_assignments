package hubs

import (
	"errors"
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
