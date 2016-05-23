package server

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/bobinette/deadpool/pingpong/protos"
)

// ----------------------------------------------------------------------------
// Sound
type Sound int32

const (
	Ping    Sound = iota
	Pong    Sound = iota
	Silence Sound = iota
)

func (s Sound) String() string {
	switch s {
	case Ping:
		return "Ping"
	case Pong:
		return "Pong"
	case Silence:
		return "Silence"
	}
	return "Unknown"
}

// ----------------------------------------------------------------------------
// Player

type Player struct {
	Id    int32
	Sound Sound
}

// ----------------------------------------------------------------------------
// Game
type Game interface {
	Play(Sound) (*protos.PlayReply, error)
	History() []*protos.GameEvent
}

type game struct {
	history []*protos.GameEvent
}

func NewGame() Game {
	return &game{
		history: make([]*protos.GameEvent, 0),
	}
}

func (g *game) Play(s Sound) (*protos.PlayReply, error) {
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}

	e := &protos.GameEvent{
		Sound:     int32(s),
		Timestamp: ts,
	}
	g.history = append(g.history, e)

	rep := protos.PlayReply{
		Accepted: true,
	}
	return &rep, nil
}

func (g *game) History() []*protos.GameEvent {
	return g.history
}
