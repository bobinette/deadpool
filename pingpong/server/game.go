package server

import (
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/bobinette/deadpool/pingpong/protos"
)

type Game interface {
	Play(int32) (*protos.PlayReply, error)
	Status() *protos.GameStatus
}

type game struct {
	history []*protos.GameEvent
}

func NewGame() Game {
	return &game{
		history: make([]*protos.GameEvent, 0),
	}
}

func (g *game) Play(id int32) (*protos.PlayReply, error) {
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, err
	}

	e := &protos.GameEvent{
		PlayerId:  id,
		Timestamp: ts,
	}
	g.history = append(g.history, e)

	rep := protos.PlayReply{
		Accepted: true,
	}
	return &rep, nil
}

func (g *game) Status() *protos.GameStatus {
	return &protos.GameStatus{
		CurrentPlayer: -1,
		Pingpong:      g.history,
	}
}
