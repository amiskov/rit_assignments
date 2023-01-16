package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"

	"ritsockets/console"
	"ritsockets/handler"
	"ritsockets/hubs"
)

var hubSize = flag.Int("s", 3, "hub size")

func main() {
	flag.Parse()

	tmpl := template.Must(template.ParseFiles("index.html"))
	wsHubs := hubs.NewHubsStorage(*hubSize)
	wsHandler := handler.NewWsHandler(wsHubs)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})
	http.Handle("/ws", wsHandler)

	go http.ListenAndServe(":8080", nil)
	fmt.Println("Server started at http://localhost:8080")
	console.New(wsHubs).Run()
}
