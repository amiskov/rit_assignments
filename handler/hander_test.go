package handler

import (
	"net/http/httptest"
	"ritsockets/hubs"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestBroadcast(t *testing.T) {
	suite := []struct {
		name          string
		msg           string
		hubSize       int
		wantHubsCount int
		clientsCount  int
	}{
		{
			name:          "successfully broadcasted to the one hub filled with clients",
			msg:           "Broadcast testing.",
			hubSize:       3,
			wantHubsCount: 1,
			clientsCount:  3,
		},
		{
			name:          "successfully broadcasted to the 2 hubs filled with clients",
			msg:           "Broadcast testing.",
			hubSize:       3,
			wantHubsCount: 2,
			clientsCount:  5,
		},
	}

	for _, tt := range suite {
		t.Run(tt.name, func(t *testing.T) {
			db := hubs.NewHubsDB(tt.hubSize)
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
