package console

import (
	"fmt"
	"ritsockets/hubs"
	"strings"
)

type console struct {
	hdb *hubs.HubsDB
}

func New(hdb *hubs.HubsDB) *console {
	return &console{
		hdb: hdb,
	}
}

func (c *console) Run() {
	fmt.Println("Enter command or 'exit':")
	for {
		fmt.Printf("> ")
		var cmd, param string
		fmt.Scanf("%s %s", &cmd, &param)

		switch cmd {
		case "help":
			fmt.Println("Usage:\n send --hub=<hub id>\n sendc --id=<client id>\n exit")
		case "sendc":
			clientIdParam, err := parseCommandParam(param)
			if err != nil {
				fmt.Println("bad client id:", param)
				continue
			}

			client, err := c.hdb.GetClientById(clientIdParam)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = client.SendMessage(fmt.Sprintf("Hello, client #%s!", client))
			if err != nil {
				fmt.Printf("failed sending message to the client %q; %s\n", client, err)
				continue
			}

			fmt.Println("Message sent to the client", client)
		case "send":
			hubIdParam, err := parseCommandParam(param)
			if err != nil {
				fmt.Println("bad hub id:", param)
				continue
			}

			hub, err := c.hdb.GetHubById(hubIdParam)
			if err != nil {
				fmt.Println("Failed to get Hub by ID: ", err)
				continue
			}

			hub.Broadcast(fmt.Sprintf("broadcasting to the hub %q", param))
		case "exit":
			fmt.Println("exiting...")
		default:
			fmt.Printf("Unknown command %q\n", cmd)
		}
	}
}

func parseCommandParam(param string) (string, error) {
	parts := strings.Split(param, "=")
	if len(parts) < 2 {
		return "", fmt.Errorf("bad params %v", parts)
	}
	return parts[1], nil
}
