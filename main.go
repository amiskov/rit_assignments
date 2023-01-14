package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"ritsockets/console"
	"ritsockets/hubs"
)

const hubSize = 3

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	tmpl := template.Must(template.ParseFiles("index.html"))
	wsHubs := hubs.NewHubsDB(hubSize)
	cons := console.New(wsHubs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		wsHubs.Add(conn)
	})

	go http.ListenAndServe(":8080", nil)
	fmt.Println("Server started at http://localhost:8080")
	cons.Run()
}
