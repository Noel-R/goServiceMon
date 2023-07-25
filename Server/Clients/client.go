package Clients

import (
	"Generics"
	"github.com/google/uuid"
	"net"
)

const (
	CONNECTED = iota
	DISCONNECTED
	UNASSIGNED
	ASSIGNED
)

type client struct {
	UID      uuid.UUID
	state    int
	socket   net.Conn
	messages map[int]Generics.Message
}
