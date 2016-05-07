package main

import (
	"math/rand"
	"time"
)

type Player interface {
	Name() string
	Play() string
}

type player struct {
	messages []string
}

func NewPlayer() Player {
	return &player{
		messages: []string{"Pizza? Yolo!", "Hello, you", "I am Gonzalez, Speedy Gonzalez"},
	}
}

func (p *player) Name() string {
	return "Speedy Gonzalez"
}

func (p *player) Play() string {
	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)

	return p.messages[gen.Intn(len(p.messages))]
}
