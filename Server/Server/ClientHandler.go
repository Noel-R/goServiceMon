package Server

import (
	"Server/Client"
	"fmt"
	"time"
)

type ClientHandler struct {
	Client  *Client.Client
	Created time.Time
}

type IClientHandler interface {
}

func (c *ClientHandler) New(client *Client.Client) *ClientHandler {
	return &ClientHandler{client, time.Now()}
}

func (c *ClientHandler) GetClient() *Client.Client {
	return c.Client
}

func (c *ClientHandler) Start() {
	defer func(c *ClientHandler) {
		client := c.GetClient()
		switch client.GetState() {
		case Client.DISCONNECTED:

			break
		case Client.ASSIGNED:

			break
		default:
			fmt.Printf("Client in an unknown sate, %v\n, %v", client.GetID(), client.GetState())
			return
		}
		fmt.Printf("Client %v disconnected.", client.GetID().String())
	}(c)

	client := c.GetClient()

	for {
		if state := client.GetState(); !(state == Client.UNASSIGNED || state == Client.CONNECTED) {
			return
		}
	}
}
