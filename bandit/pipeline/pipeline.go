package pipeline

import (
	"fmt"
)

type Pipeline struct {
	inputChs  map[string][]*Channel
	outputChs map[string]*Channel
	mounted   bool

	// Map outputCh.Key -> Unit
	units []Unit
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		inputChs:  make(map[string][]*Channel),
		outputChs: make(map[string]*Channel),
		mounted:   false,

		units: make([]Unit, 0),
	}
}

// Add adds a unit to the set. If the unit cannot be added to the set, for example because it is
// already present in the set, an error is returned.
func (p *Pipeline) Add(unit Unit) error {
	for _, och := range unit.Outputs() {
		if _, ok := p.outputChs[och.Key]; ok {
			return fmt.Errorf("%s: output channel already added", och.Key)
		}
	}

	for _, och := range unit.Outputs() {
		p.outputChs[och.Key] = och
	}

	for _, ich := range unit.Inputs() {
		ichs, ok := p.inputChs[ich.Key]
		if !ok {
			ichs = make([]*Channel, 0)
		}

		ichs = append(ichs, ich)
		p.inputChs[ich.Key] = ichs
	}

	p.units = append(p.units, unit)

	return nil
}

func (p *Pipeline) Mount() error {
	if p.mounted {
		return fmt.Errorf("already mounterd")
	}

	for key, och := range p.outputChs {
		ichs, ok := p.inputChs[key]
		if !ok {
			continue
		}

		go p.dispatch(och, ichs)
	}

	for _, u := range p.units {
		u.Start()
	}

	p.mounted = true
	return nil
}

func (p *Pipeline) dispatch(och *Channel, ichs []*Channel) {
	for {
		i := <-och.Ch

		for _, ich := range ichs {
			ich.Ch <- i
		}
	}
}
