package systems

import (
	"log"

	"github.com/bobinette/deadpool/bandit/components"
	"github.com/bobinette/deadpool/bandit/pipeline"
)

type Player struct {
	// Inputs
	arm        *pipeline.Channel
	value      *pipeline.Channel
	stateInput *pipeline.Channel

	// Ouputs
	State chan components.PlayerState
}

func NewPlayer() *Player {
	p := Player{
		arm:        pipeline.NewChannel("arm"),
		value:      pipeline.NewChannel("value"),
		stateInput: pipeline.NewChannel("state"),

		State: make(chan components.PlayerState),
	}

	return &p
}

func (p *Player) Inputs() []*pipeline.Channel {
	return []*pipeline.Channel{p.arm, p.value, p.stateInput}
}

func (p *Player) Outputs() []*pipeline.Channel {
	return nil
}

func (p *Player) Start() {
	go p.run()
}

func (p *Player) run() {
	for {
		iarm := <-p.arm.Ch
		arm, ok := iarm.(int32)
		if !ok {
			log.Fatalf("Expected arm to be of type int32, got %T", iarm)
		}

		istate := <-p.stateInput.Ch
		state, ok := istate.(components.PlayerState)
		if !ok {
			log.Fatalf("Expected arm to be of type PlayerState, got %T", istate)
		}

		ivalue := <-p.value.Ch
		value, ok := ivalue.(float64)
		if !ok {
			log.Fatalf("Expected arm to be of type float, got %T", ivalue)
		}

		state.Knowledge[int(arm)] = value
		state.Score += value
		state.RemainingPlies -= 1
		p.State <- state
	}
}
