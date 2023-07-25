package main

import (
	"fmt"
)

func main() {

	var s iServer

	s = &server{}
	s = s.initServer("127.0.0.1:")

	fmt.Print("Starting listener.\n")

	s.createListener()

	addr, port := s.getAddress()

	fmt.Printf("Listener started on: %v:%v\n", addr, port)
	for {
		continue
	}
}
