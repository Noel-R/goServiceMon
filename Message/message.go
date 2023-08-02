package Message

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net"
	"strconv"
	"time"
)

const (
	Ready = iota
	Sent
	Errored
)

type Message struct {
	msgID  uuid.UUID
	destID uuid.UUID
	srcID  uuid.UUID

	payload map[string][]string
	time    time.Time
	msgType int

	msgBody []byte

	state int
}

type IMessage interface {
	New(uuid.UUID, uuid.UUID, []string, int) (*Message, error)
	Send(conn *net.TCPConn) error
	Body() []byte
}

func (m *Message) New(sendID uuid.UUID, destID uuid.UUID, body []string, msgType int) (*Message, error) {
	var err error
	m.msgID, m.srcID, m.destID = uuid.New(), sendID, destID
	m.msgType = msgType
	body = append(body, time.Now().String())
	body = append(body, strconv.Itoa(m.msgType))
	m.payload = map[string][]string{m.msgID.String(): body}
	m.msgBody, err = json.Marshal(m.payload)
	if err != nil {
		fmt.Printf("Error creating message: %v\n", err)
		m.state = Errored
		return nil, err
	}
	m.state = Ready
	return m, err
}

func (m *Message) Send(conn *net.TCPConn) error {
	if _, err := conn.Read(make([]byte, 1)); err != io.EOF && m.state == Ready {
		if _, err := conn.Write(m.msgBody); err != nil {
			m.state = Errored
			return err
		}
	}
	m.state = Sent
	return nil
}

func (m *Message) Body() []byte {
	return m.msgBody
}

func Decode(received []byte) string {
	decoded := string(received)
	return decoded
}

func (m *Message) GetID() uuid.UUID {
	return m.msgID
}
