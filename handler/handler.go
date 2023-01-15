package handler

import (
	"fmt"
	"log"
	"net/http"

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
	Add(*websocket.Conn)
}

type hubsHandler struct {
	db DB
}

func New(db DB) *hubsHandler {
	return &hubsHandler{
		db: db,
	}
}

func (h hubsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		fmt.Println("defer is called...")
		// conn.Close()
		// fmt.Println("...connection is closed.")
	}()

	h.db.Add(conn)
}
