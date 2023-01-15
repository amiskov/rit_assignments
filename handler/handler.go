package handler

import (
	"encoding/json"
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

func New(db DB) *hubsHandler {
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

	// TODO: the body of `for` loop is actually don't used for this task.

	for {
		_, m, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		var in inbound
		err = json.Unmarshal(m, &in)
		if err != nil {
			// handleError(ws, err)
			log.Printf("error: %v", err)
			continue
		}

		out, err := json.Marshal(outbound{Body: "Body of outbound message"})
		if err != nil {
			// handleError(ws, err)
			log.Printf("error: %v", err)
			continue
		}

		err = ws.WriteMessage(websocket.BinaryMessage, out)
		if err != nil {
			// handleError(ws, err)
			log.Printf("error: %v", err)
			continue
		}
	}
}
