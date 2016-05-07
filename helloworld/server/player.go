package helloworld

import (
	"time"
)

type Player struct {
	Id   int
	Name string

	start time.Time
}

func NewPlayer(id int, name string) *Player {
	return &Player{
		Id:    id,
		Name:  name,
		start: time.Now(),
	}
}

func (p *Player) ConnectedFor() time.Duration {
	return time.Since(p.start)
}
