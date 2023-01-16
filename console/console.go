package console

import (
	"fmt"
	"log"

	"ritsockets/hubs"
)

type Storage interface {
	GetHubById(string) (*hubs.Hub, error)
	GetClientById(string) (*hubs.Client, error)
}

type console struct {
	db Storage
}

func New(db Storage) *console {
	return &console{
		db: db,
	}
}

func (c *console) Run() {
	fmt.Println("Enter command:")
	for {
		fmt.Printf("> ")
		var cmd, id string
		fmt.Scanf("%s %s", &cmd, &id)

		switch cmd {
		case "help":
			fmt.Println("Usage:\n send HUB_ID\n sendc CLIENT_ID")
		case "sendc":
			client, err := c.db.GetClientById(id)
			if err != nil {
				log.Println(err)
				continue
			}

			err = client.SendMessage(fmt.Sprintf("Hello, client #%s!", client))
			if err != nil {
				log.Printf("failed sending message to the client %q; %s\n", client, err)
				continue
			}

			fmt.Println("Message sent to the client", client)
		case "send":
			hub, err := c.db.GetHubById(id)
			if err != nil {
				log.Println("Failed to get Hub by ID: ", err)
				continue
			}

			hub.Broadcast(fmt.Sprintf("broadcasting to the hub %s", id))
		default:
			fmt.Printf("Unknown command %q\n", cmd)
		}
	}
}
