package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icrowley/fake"

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

	hdb := hubs.NewHubsDB(hubSize)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		hdb.Add(conn)

		// go sendNewMsgNotifications(conn)
	})

	fmt.Println("starting server at :8080")
	go http.ListenAndServe(":8080", nil)
	console()
}

func console() {
	fmt.Println("Enter command or 'exit':")
	for {
		fmt.Printf("> ")
		var cmd string
		fmt.Scanln(&cmd)

		if strings.HasPrefix(cmd, "sendc ") {
			fmt.Println("sending to the client")
			continue
		}

		if strings.HasPrefix(cmd, "send ") {
			fmt.Println("broadcasting to hub")
			continue
		}

		if cmd == "exit" {
			fmt.Println("exiting...")
			continue
		}

		fmt.Println("Unknown command. Use `send --hub=<hub id>`, `sendc --id=<client id>` or `exit`.")
	}
}

func sendNewMsgNotifications(client *websocket.Conn) {
	ticker := time.NewTicker(2 * time.Second)
	for {
		w, err := client.NextWriter(websocket.TextMessage)
		if err != nil {
			ticker.Stop()
			break
		}

		msg := newMessage()
		w.Write(msg)
		w.Close()

		<-ticker.C
	}
}

func newMessage() []byte {
	data, _ := json.Marshal(map[string]string{
		"email":   fake.EmailAddress(),
		"name":    fake.FirstName() + " " + fake.LastName(),
		"subject": fake.Product() + " " + fake.Model(),
	})
	return data
}