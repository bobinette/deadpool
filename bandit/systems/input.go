package systems

import (
	// "log"

	"github.com/bobinette/deadpool/bandit/components"
	"github.com/bobinette/deadpool/bandit/pipeline"
)

type Input struct {
	Arm         chan int32
	Bandit      chan components.Arms
	PlayerState chan components.PlayerState

	arm         *pipeline.Channel
	bandit      *pipeline.Channel
	playerState *pipeline.Channel
}

func NewInput() *Input {
	return &Input{
		Arm:         make(chan int32),
		Bandit:      make(chan components.Arms),
		PlayerState: make(chan components.PlayerState),

		arm:         pipeline.NewChannel("arm"),
		bandit:      pipeline.NewChannel("bandit"),
		playerState: pipeline.NewChannel("state"),
	}
}

func (i *Input) Inputs() []*pipeline.Channel { return nil }

func (i *Input) Outputs() []*pipeline.Channel {
	return []*pipeline.Channel{i.arm, i.bandit, i.playerState}
}

func (i *Input) Start() {
	go func() {
		for {
			b := <-i.Bandit
			i.bandit.Ch <- b
		}
	}()

	go func() {
		for {
			s := <-i.PlayerState
			i.playerState.Ch <- s
		}
	}()

	go func() {
		for {
			a := <-i.Arm
			i.arm.Ch <- a
		}
	}()
}
