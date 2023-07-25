package Generics

import (
	"github.com/google/uuid"
	"io"
	"net"
	"strings"
	"time"
)

const (
	Ready = iota
	Sent
	Errored
)

type Message struct {
	destID uuid.UUID
	srcID  uuid.UUID
	body   []string
	time   time.Time
	state  int
}

func (m Message) createMsg(dA uuid.UUID, sA uuid.UUID, contents ...string) Message {
	m.destID, m.srcID, m.body, m.time, m.state = dA, sA, contents, time.Now(), Ready
	return m
}

func (m Message) sendMsg(s net.TCPConn) (bool, *Message) {
	if _, err := s.Read(make([]byte, 1)); err != io.EOF && m.state == Ready {
		if _, err := s.Write([]byte(strings.Join(m.body, "::"))); err != nil {
			m.state = Errored
			return false, &m
		}
		m.state = Sent
	}
	return true, &m
}
