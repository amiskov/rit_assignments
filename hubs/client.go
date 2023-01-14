package hubs

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type clientID uuid.UUID

type client struct {
	id   clientID
	conn *websocket.Conn
}

func (c *client) String() string {
	return uuid.UUID(c.id).String()
}
