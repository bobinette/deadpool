package bandit

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bobinette/deadpool/bandit/components"
	"github.com/bobinette/deadpool/bandit/pipeline"
	"github.com/bobinette/deadpool/bandit/systems"
)

type GameManager struct {
	bandit      components.Arms        // Bandit components
	playerState components.PlayerState // Player components

	InputChannel chan components.InputEvent

	input     *systems.Input
	extractor *systems.Extractor
	player    *systems.Player
	pipeline  *pipeline.Pipeline
}

func NewGameManager(n int32) (*GameManager, error) {
	gm := GameManager{
		bandit: make(components.Arms, n),
		playerState: components.PlayerState{
			Knowledge:      make(components.Arms, int(n)),
			RemainingPlies: n,
			Score:          0,
		},

		InputChannel: make(chan components.InputEvent),

		input:     systems.NewInput(),
		extractor: systems.NewExtractor(),
		player:    systems.NewPlayer(),
		pipeline:  pipeline.NewPipeline(),
	}

	// ---- Pipe units
	if err := gm.pipeline.Add(gm.input); err != nil {
		return nil, fmt.Errorf("could not register input: %v", err)
	}
	if err := gm.pipeline.Add(gm.extractor); err != nil {
		return nil, fmt.Errorf("could not register extractor: %v", err)
	}
	if err := gm.pipeline.Add(gm.player); err != nil {
		return nil, fmt.Errorf("could not register player: %v", err)
	}
	if err := gm.pipeline.Mount(); err != nil {
		return nil, fmt.Errorf("could not mount pipeline: %v", err)
	}

	return &gm, nil
}

func (gm *GameManager) Start() {
	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)
	for i := range gm.bandit {
		gm.bandit[i] = gen.NormFloat64()
	}
	log.Println(gm.bandit)

	gm.pipeline.Mount()
	go gm.run()
}

func (gm *GameManager) run() {
	for {
		evt, ok := <-gm.InputChannel
		if !ok {
			return
		}

		gm.input.Arm <- evt.Input.Arm
		gm.input.Bandit <- gm.bandit
		gm.input.PlayerState <- gm.playerState

		s := <-gm.player.State
		gm.playerState = s
		evt.Output <- s
	}
}
