package systems

import (
	"log"

	"github.com/bobinette/deadpool/bandit/components"
	"github.com/bobinette/deadpool/bandit/pipeline"
)

type Extractor struct {
	// Inputs
	arm    *pipeline.Channel
	bandit *pipeline.Channel

	// Ouputs
	value *pipeline.Channel
}

func NewExtractor() *Extractor {
	e := Extractor{
		arm:    pipeline.NewChannel("arm"),
		bandit: pipeline.NewChannel("bandit"),

		value: pipeline.NewChannel("value"),
	}

	return &e
}

func (e *Extractor) Inputs() []*pipeline.Channel {
	return []*pipeline.Channel{e.arm, e.bandit}
}

func (e *Extractor) Outputs() []*pipeline.Channel {
	return []*pipeline.Channel{e.value}
}

func (e *Extractor) Start() {
	go e.run()
}

func (e *Extractor) run() {
	for {
		iarm := <-e.arm.Ch
		arm, ok := iarm.(int32)
		if !ok {
			log.Fatalf("Expected arm to be of type int32, got %T", iarm)
		}

		ibandit := <-e.bandit.Ch
		bandit, ok := ibandit.(components.Arms)
		if !ok {
			log.Fatalf("Expected arm to be of type Arms, got %T", ibandit)
		}

		e.value.Ch <- bandit[int(arm)]
	}
}
