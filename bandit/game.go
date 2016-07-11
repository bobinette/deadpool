package bandit

import (
	"log"
	"math/rand"
	"time"
)

type InputEvent struct {
	Input  Input
	Output chan float64
}

type IO struct {
	Input  Input
	Output float64
}

// Game is the entity component manager, also responsible for the game main loop
type Game struct {
	Player Player
	Bandit Bandit

	IO           IO
	InputChannel chan InputEvent

	systems []System
}

func NewGame() *Game {
	g := Game{
		InputChannel: make(chan InputEvent),
	}

	g.systems = []System{
		&BanditSystem{},
		&PlayerSystem{},
	}

	for _, s := range g.systems {
		s.Init(&g)
	}
	return &g
}

func (g *Game) Start() {
	for {
		// Wait for input
		e, ok := <-g.InputChannel
		if !ok {
			return
		}

		// Deal with input
		g.IO = IO{
			Input:  e.Input,
			Output: 0,
		}

		// Update game
		for _, s := range g.systems {
			s.Update(g)
		}

		// Manager ouput
		// @TODO: handle notifications, handler more elegantly
		// the ouput workflow
		e.Output <- g.IO.Output
	}
}

// ---- System ----
type System interface {
	Init(*Game)
	Update(*Game)
}

// ---- Bandit system ----
type BanditSystem struct{}

func (s *BanditSystem) Init(g *Game) {
	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)

	b := &g.Bandit
	b.Values = make([]float64, 10)
	for i := range b.Values {
		b.Values[i] = gen.NormFloat64()
	}

	log.Printf("%+v", b.Values)
}

func (s *BanditSystem) Update(g *Game) {}

// ---- Player system ----
type PlayerSystem struct{}

func (s *PlayerSystem) Init(g *Game) {
	p := &g.Player
	b := g.Bandit

	p.History = make([]int32, len(b.Values))
	p.Score = 0
	p.RemainingPlies = int32(len(b.Values))
	p.Knowledge = make([]float64, len(b.Values))
}

func (s *PlayerSystem) Update(g *Game) {
	p := &g.Player
	if p.RemainingPlies <= 0 {
		return
	}

	b := g.Bandit
	a := g.IO.Input.Arm

	if a < 0 || a >= int32(len(b.Values)) {
		return
	}

	v := b.Values[a]
	p.Score += v
	p.History[a] -= 1
	p.RemainingPlies -= 1

	g.IO.Output = v
}
