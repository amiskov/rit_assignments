package handler

import (
	"log"
	"net/http"
	"ritsockets/hubs"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // we don't care about CORS here.
	},
}

type DB interface {
	Add(*hubs.Client)
}

type hubsHandler struct {
	db DB
}

func NewWsHandler(db DB) *hubsHandler {
	return &hubsHandler{
		db: db,
	}
}

type inbound struct{}
type outbound struct {
	Body string `json:"body"`
}

func (h hubsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	c := hubs.NewClient(ws)

	defer func() {
		ws.Close()
		log.Printf("connection is closed for client %q\n", c)
	}()

	h.db.Add(c)

	// Dead simple echo server which sends the received message back to the client.
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}
