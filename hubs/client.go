package hubs

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id   uuid.UUID
	conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Id:   uuid.New(),
		conn: conn,
	}
}

func (c *Client) String() string {
	return uuid.UUID(c.Id).String()
}

func (c *Client) SendMessage(msg string) error {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	w.Write([]byte(msg))
	w.Close()
	return nil
}
