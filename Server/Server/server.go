package Server

import (
	"Server/Client"
	"fmt"
	"github.com/google/uuid"
	"net"
)

type IServer interface {
	InitServer(string) *Server
	CreateListener()
	GetAddress() (string, string)
	Socket() *net.TCPListener
	SetSocket(*net.TCPListener)
	Connections() map[int]*net.TCPConn
	SetConnections(map[int]*net.TCPConn)
	Address() string
	SetAddress(string)
	GetID() uuid.UUID
}

type Server struct {
	UID         uuid.UUID
	socket      *net.TCPListener
	connections map[int]*net.TCPConn
	address     string
	Clients     map[uuid.UUID]*Client.Client
}

func (s *Server) Socket() *net.TCPListener {
	return s.socket
}

func (s *Server) SetSocket(socket *net.TCPListener) {
	s.socket = socket
}

func (s *Server) Connections() map[int]*net.TCPConn {
	return s.connections
}

func (s *Server) SetConnections(connections map[int]*net.TCPConn) {
	s.connections = connections
}

func (s *Server) Address() string {
	return s.address
}

func (s *Server) SetAddress(address string) {
	s.address = address
}

func (s *Server) InitServer(addr string) *Server {
	return &Server{
		UID:         uuid.New(),
		connections: make(map[int]*net.TCPConn),
		address:     addr,
	}
}

func (s *Server) GetAddress() (string, string) {
	addr, port, err := net.SplitHostPort(s.socket.Addr().String())
	if err != nil {
		fmt.Printf("Error during address retrieval: %v", err)
		return "", ""
	}
	return addr, port
}

func (s *Server) CreateListener() {
	s.address = "127.0.0.1:"
	Addr, err := net.ResolveTCPAddr("tcp", s.address)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", Addr)
	if err != nil {
		panic(err)
	}
	s.socket = listener
	go s.StartListener()
}

func (s *Server) StartListener() {
	defer func(s *Server) {
		err := s.socket.Close()
		if err != nil {
			panic(err)
		}
	}(s)
	for {
		conn, err := s.socket.AcceptTCP()
		if err != nil {
			panic(err)
		}
		s.AddConn(conn)
		go ConnectionHandler(s, conn)
	}
}

func (s *Server) AddConn(conn *net.TCPConn) {
	s.connections[len(s.connections)+1] = conn
	return
}

func (s *Server) RmvConn(conn *net.TCPConn) {
	for k, v := range s.connections {
		if v == conn {
			delete(s.connections, k)
			return
		}
	}
}

func (s *Server) GetID() uuid.UUID {
	return s.UID
}

func ConnectionHandler(s *Server, conn *net.TCPConn) {
	defer func(conn *net.TCPConn) {
		err := conn.Close()
		fmt.Printf("Connection closed.\n")
		if err != nil {
			fmt.Printf("Conn returned an error during closure: %#v\n", err)
		}
		// Delete from connections map.
		s.RmvConn(conn)
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
			fmt.Printf("Conn encountered an error during send: %#v\n", err)
			return
		}
	}
}
