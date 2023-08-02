package Client

import (
	"Message"
	"Message/Type"
	"Server/Server"
	"fmt"
	"github.com/google/uuid"
	"net"
)

const (
	CONNECTED = iota
	DISCONNECTED
	UNASSIGNED
	ASSIGNED
)

type Client struct {
	UID      uuid.UUID
	state    int
	socket   *net.TCPConn
	messages map[int]*Message.Message
}

type IClient interface {
	New(*net.TCPConn) *Client
	SetState(int)
	GetState() int
	AppendMessage(*Message.Message)
	RemoveMessage(Message.Message)
	GetID() uuid.UUID
	GetMessages() *map[int]*Message.Message
	SendMessage(*Message.Message)
}

func (c *Client) SetConnection(conn *net.TCPConn) {
	c.socket = conn
	return
}

func (c *Client) GenerateID() {
	c.UID = uuid.New()
	return
}

func (c *Client) SetState(state int) {
	c.state = state
	return
}

func (c *Client) GetState() int {
	return c.state
}

func (c *Client) AppendMessage(m *Message.Message) {
	c.messages[len(c.messages)+1] = m
	return
}

func (c *Client) RemoveMessage(m Message.Message) {
	for i, msg := range c.messages {
		if msg.GetID() == m.GetID() {
			delete(c.messages, i)
			return
		}
	}
}

func (c *Client) New(s *net.TCPConn) *Client {
	c.SetConnection(s)
	c.GenerateID()
	c.messages = map[int]*Message.Message{}
	c.SetState(CONNECTED)
	return c
}

func (c *Client) GetID() uuid.UUID {
	return c.UID
}

func (c *Client) GetMessages() *map[int]*Message.Message {
	return &c.messages
}

func (c *Client) SendMessage(m *Message.Message) error {
	err := m.Send(c.socket)
	if err != nil {
		return err
	}
	c.AppendMessage(m)
	return err
}

func (c *Client) Disconnect(s *Server.Server) {
	m := Message.Message{}
	message, _ := m.New(s.GetID(), c.GetID(), []string{"Disconnected"}, Type.DISCONNECT)
	if err := c.SendMessage(message); err != nil {
		fmt.Printf("Client disconnected prematurely, %v.", c.GetID())
		c.SetState(DISCONNECTED)
		return
	}
	fmt.Printf("Client disconnected, %v.", c.GetID())
	c.SetState(DISCONNECTED)
	return
}
