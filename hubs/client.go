package hubs

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type clientID uuid.UUID

type Client struct {
	id   clientID
	conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		id:   clientID(uuid.New()),
		conn: conn,
	}
}

func (c *Client) String() string {
	return uuid.UUID(c.id).String()
}

func (c *Client) SendMessage(msg string) error {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(map[string]string{
		"message": msg,
	})

	w.Write(data)
	w.Close()
	return nil
}
