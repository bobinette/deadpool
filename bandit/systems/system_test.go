package systems

import (
	"reflect"
	"testing"
	"time"

	"github.com/bobinette/deadpool/bandit/components"
	"github.com/bobinette/deadpool/bandit/pipeline"
)

func TestSystem(t *testing.T) {
	input := NewInput()
	extractor := NewExtractor()
	player := NewPlayer()

	var tts = map[string]struct {
		bandit   components.Arms
		state    components.PlayerState
		arm      int32
		expected components.PlayerState
	}{
		"one arm bandit": {
			bandit: []float64{2},
			state: components.PlayerState{
				Knowledge:      []float64{0},
				RemainingPlies: 1,
				Score:          0,
			},
			arm: 0,
			expected: components.PlayerState{
				Knowledge:      []float64{2},
				RemainingPlies: 0,
				Score:          2,
			},
		},
		"two arm bandit": {
			bandit: []float64{3, 2},
			state: components.PlayerState{
				Knowledge:      []float64{3, 0},
				RemainingPlies: 1,
				Score:          3,
			},
			arm: 1,
			expected: components.PlayerState{
				Knowledge:      []float64{3, 2},
				RemainingPlies: 0,
				Score:          5,
			},
		},
	}

	p := pipeline.NewPipeline()
	if err := p.Add(input); err != nil {
		t.Fatalf("could not register input: %v", err)
	}
	if err := p.Add(extractor); err != nil {
		t.Fatalf("could not register extractor: %v", err)
	}
	if err := p.Add(player); err != nil {
		t.Fatalf("could not register player: %v", err)
	}
	if err := p.Mount(); err != nil {
		t.Fatalf("could not mount pipeline: %v", err)
	}

	d := 100 * time.Millisecond
	for name, tt := range tts {
		pushed := false
		select {
		case input.Arm <- tt.arm:
			pushed = true
		case <-time.After(d):
			t.Errorf("%s - Could not send arm, waited for %v", name, d)
		}
		if !pushed {
			continue
		}

		pushed = false
		select {
		case input.Bandit <- tt.bandit:
			pushed = true
		case <-time.After(d):
			t.Errorf("%s - Could not send bandit, waited for %v", name, d)
		}
		if !pushed {
			continue
		}

		pushed = false
		select {
		case input.PlayerState <- tt.state:
			pushed = true
		case <-time.After(d):
			t.Errorf("%s - Could not send state, waited for %v", name, d)
		}
		if !pushed {
			continue
		}

		select {
		case state := <-player.State:
			if !reflect.DeepEqual(state, tt.expected) {
				t.Errorf("%s - Incorrect state: expected %v got %v", name, tt.expected, state)
			}
		case <-time.After(d):
			t.Errorf("%s - Did not receive state, waited for %v", name, d)
		}
	}
}
