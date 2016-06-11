package pingpong

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/bobinette/deadpool/pingpong/proto"
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
	Name  string
}

// ----------------------------------------------------------------------------
// Game
type Game interface {
	Play(Sound) (*proto.PlayReply, error)
	History() []*proto.GameEvent
}

type game struct {
	history []*proto.GameEvent
}

func NewGame() Game {
	return &game{
		history: make([]*proto.GameEvent, 0),
	}
}

func (g *game) Play(s Sound) (*proto.PlayReply, error) {
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}

	e := &proto.GameEvent{
		Sound:     int32(s),
		Timestamp: ts,
	}
	g.history = append(g.history, e)

	rep := proto.PlayReply{
		Accepted: true,
	}
	return &rep, nil
}

func (g *game) History() []*proto.GameEvent {
	return g.history
}
