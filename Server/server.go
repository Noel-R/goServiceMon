package main

import (
	"fmt"
	"github.com/google/uuid"
	"net"
)

type iServer interface {
	initServer(string) *server
	createListener()
	getAddress() (string, string)
	Socket() net.Listener
	SetSocket(net.Listener)
	Connections() map[int]net.Conn
	SetConnections(map[int]net.Conn)
	Address() string
	SetAddress(string)
}

type server struct {
	UID         uuid.UUID
	socket      net.Listener
	connections map[int]net.Conn
	address     string
}

func (s *server) Socket() net.Listener {
	return s.socket
}

func (s *server) SetSocket(socket net.Listener) {
	s.socket = socket
}

func (s *server) Connections() map[int]net.Conn {
	return s.connections
}

func (s *server) SetConnections(connections map[int]net.Conn) {
	s.connections = connections
}

func (s *server) Address() string {
	return s.address
}

func (s *server) SetAddress(address string) {
	s.address = address
}

func (s *server) initServer(addr string) *server {
	return &server{
		UID:         uuid.New(),
		socket:      nil,
		connections: make(map[int]net.Conn),
		address:     addr,
	}
}

func (s *server) getAddress() (string, string) {
	addr, port, err := net.SplitHostPort(s.socket.Addr().String())
	if err != nil {
		fmt.Printf("Error during address retrieval: %v", err)
		return "", ""
	}
	return addr, port
}

func (s *server) createListener() {
	s.address = "127.0.0.1:"

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		panic(err)
	}

	s.socket = listener
	go s.startListener()
}

func (s *server) startListener() {
	defer func(s *server) {
		err := s.socket.Close()
		if err != nil {
			panic(err)
		}
	}(s)
	for {
		conn, err := s.socket.Accept()
		if err != nil {
			panic(err)
		}
		s.connections[len(s.connections)+1] = conn
		go connectionHandler(s, conn)
	}
}

func connectionHandler(s *server, conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		fmt.Printf("Connection closed.")
		if err != nil {
			fmt.Printf("Conn returned an error during closure: %#v\n", err)
			// Delete from connections map.
			for k, v := range s.connections {
				if v == conn {
					delete(s.connections, k)
					return
				}
			}
		}
	}(conn)
	for {
		buffer := make([]byte, 1024)

		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Conn encountered an error during messaging: %#v\n", err)
			return
		}

		fmt.Printf("Received: %v\n", string(buffer))
		_, err = conn.Write([]byte("Message Received."))
		if err != nil {
			fmt.Printf("Conn encountered an error during messaging: %#v\n", err)
			return
		}
	}
}
