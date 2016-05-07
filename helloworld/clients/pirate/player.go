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
		messages: []string{
			"Prepare to be boarded!",
			"Surrrrrender the booty!",
			"Save a ship. Ride a pirate.",
			"Take what you can, give nothing back",
			"If ye can’t trust a pirate, ye damn well can’t trust a merchant either!",
		},
	}
}

func (p *player) Name() string {
	return "Jack Sparrow"
}

func (p *player) Play() string {
	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)

	return p.messages[gen.Intn(len(p.messages))]
}
