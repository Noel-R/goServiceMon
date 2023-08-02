package main

import (
	"Message"
	"Message/Type/Register"
	"Server/Server"
	"fmt"
	"github.com/google/uuid"
)

func Start() {
	var s Server.IServer

	s = &Server.Server{}
	s = s.InitServer("127.0.0.1:")

	fmt.Print("Starting listener.\n")

	s.CreateListener()

	addr, port := s.GetAddress()

	fmt.Printf("Listener started on: %v:%v\n", addr, port)
}

func main() {
	// Starts server.
	//
	Start()

	//Example message creation code.

	msgArray := make(map[int]*Message.Message)
	for i := 0; i <= 10; i++ {
		message := Message.Message{}
		if m, err := message.New(uuid.New(), uuid.New(), []string{"test", "message"}, Register.Request); err != nil {
			fmt.Printf("Error: %v", err)
		} else {
			msgArray[i] = m
		}
	}
	for k, v := range msgArray {
		fmt.Printf("%v: %v\n", k, Message.Decode(v.Body()))
	}

	for {
		continue
	}
}
