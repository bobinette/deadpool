package helloworld

import (
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/bobinette/deadpool/helloworld/protos"
)

type Game interface {
	Players() int
	Play(msg string, p *Player) (*protos.GameEvent, error)

	Status() *protos.GameStatus
	Reset()
}

type game struct {
	messages []message
}

func NewGame() Game {
	return &game{
		messages: make([]message, 0),
	}
}

func (g *game) Players() int {
	return 2
}

func (g *game) Play(msg string, p *Player) (*protos.GameEvent, error) {
	m := newMessage(msg, p.Name)

	n, err := m.ToGameEvent()
	if err != nil {
		return nil, err
	}

	g.messages = append(g.messages, m)
	return n, nil
}

func (g *game) Status() *protos.GameStatus {
	var msgs []*protos.GameEvent
	for _, msg := range g.messages {
		n, err := msg.ToGameEvent()
		if err != nil {
			log.Println(err)
			continue
		}
		msgs = append(msgs, n)
	}
	return &protos.GameStatus{
		Chat: msgs,
	}
}

func (g *game) Reset() {
	g.messages = make([]message, 0)
}

type message struct {
	text   string
	author string
	time   time.Time
}

func newMessage(msg, author string) message {
	return message{
		text:   msg,
		author: author,
		time:   time.Now(),
	}
}

func (m message) ToGameEvent() (*protos.GameEvent, error) {
	ts, err := ptypes.TimestampProto(m.time)
	if err != nil {
		return nil, err
	}
	return &protos.GameEvent{
		Message:   m.text,
		Author:    m.author,
		Timestamp: ts,
	}, nil
}
