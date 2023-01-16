package handler

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"

	"ritsockets/hubs"
)

func TestHubs(t *testing.T) {
	suite := []struct {
		name          string
		msg           string
		hubSize       int
		wantHubsCount int
		clientsCount  int
	}{
		{
			name:          "create 1 hub of size 3 with 3 clients",
			hubSize:       3,
			clientsCount:  3,
			wantHubsCount: 1,
		},
		{
			name:          "create 2 hubs of size 3 with 5 clients",
			hubSize:       3,
			clientsCount:  5,
			wantHubsCount: 2,
		},
		{
			name:          "create 1 empty hub",
			hubSize:       3,
			clientsCount:  0,
			wantHubsCount: 1,
		},
		{
			name:          "create 12 hubs of size 2 with 23 clients",
			hubSize:       2,
			clientsCount:  23,
			wantHubsCount: 12,
		},
	}

	for _, tt := range suite {
		t.Run(tt.name, func(t *testing.T) {
			db := hubs.NewHubsStorage(tt.hubSize)
			wsHandler := NewWsHandler(db)

			s := httptest.NewServer(wsHandler)
			defer s.Close()

			_, closeClients, err := spawnClients(tt.clientsCount, s.URL)
			if err != nil {
				t.Fatalf("failed spawning clients (%v)", err)
			}
			defer closeClients()

			createdHubs := db.ListHubs()
			if len(createdHubs) != tt.wantHubsCount {
				t.Fatalf("should create %d hubs, but created %d", tt.wantHubsCount, len(createdHubs))
			}

			createdClients := db.ListAllClients()
			if len(createdClients) != tt.clientsCount {
				t.Fatalf("should create %d clients, but created %d", tt.clientsCount, len(createdHubs))
			}
		})
	}
}

func TestBroadcast(t *testing.T) {
	suite := []struct {
		name         string
		msg          string
		hubSize      int
		clientsCount int
	}{
		{
			name:         "broadcast to all clients of a hub",
			msg:          "Broadcast testing.",
			hubSize:      3,
			clientsCount: 3,
		},
	}

	for _, tt := range suite {
		t.Run(tt.name, func(t *testing.T) {
			db := hubs.NewHubsStorage(tt.hubSize)
			wsHandler := NewWsHandler(db)

			s := httptest.NewServer(wsHandler)
			defer s.Close()

			wsClients, closeClients, err := spawnClients(tt.clientsCount, s.URL)
			if err != nil {
				t.Fatalf("failed spawning clients (%v)", err)
			}
			defer closeClients()

			createdHubs := db.ListHubs()
			if len(createdHubs) != 1 {
				t.Fatalf("should create %d hubs, but created %d", 1, len(createdHubs))
			}

			createdClients := createdHubs[0].ListClients()
			if len(createdClients) != tt.hubSize {
				t.Fatalf("should create %d clients, but created %d", tt.hubSize, len(createdHubs))
			}

			createdHubs[0].Broadcast(tt.msg)

			// All clients should receive the same message
			for _, c := range wsClients {
				_, p, err := c.ReadMessage()
				if err != nil {
					t.Fatalf("%v", err)
				}

				if string(p) != tt.msg {
					t.Fatalf("wrong message %q: %v", string(p), err)
				}
			}
		})
	}
}
func TestSendToClient(t *testing.T) {
	suite := []struct {
		name    string
		msg     string
		hubSize int
	}{
		{
			name:    "message sent to a client",
			msg:     "Message for a client.",
			hubSize: 3,
		},
	}

	for _, tt := range suite {
		t.Run(tt.name, func(t *testing.T) {
			db := hubs.NewHubsStorage(tt.hubSize)
			wsHandler := NewWsHandler(db)

			s := httptest.NewServer(wsHandler)
			defer s.Close()

			// Spawn 1 client
			wsClients, closeClients, err := spawnClients(1, s.URL)
			if err != nil {
				t.Fatalf("failed spawning clients (%v)", err)
			}
			defer closeClients()

			// Send a message to the spawned client
			client := db.ListHubs()[0].ListClients()[0]
			err = client.SendMessage(tt.msg)
			if err != nil {
				t.Fatalf("failed sending message to the client (%v)", err)
			}

			// Received message from the WS echo server should be the same as sent.
			_, p, err := wsClients[0].ReadMessage()
			if err != nil {
				t.Fatalf("%v", err)
			}
			if string(p) != tt.msg {
				t.Fatalf("wrong message %q: %v", string(p), err)
			}
		})
	}
}

// Creates a slice with `n` WebSocket clients connected to the given `url`.
func spawnClients(n int, url string) ([]*websocket.Conn, func(), error) {
	var wsClients []*websocket.Conn

	wsURL := "ws" + strings.TrimPrefix(url, "http")
	for i := 0; i < n; i++ {
		ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			return nil, nil, err
		}
		wsClients = append(wsClients, ws)
	}
	return wsClients, func() {
		for _, c := range wsClients {
			c.Close()
		}
	}, nil
}
