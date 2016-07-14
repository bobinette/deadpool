package systems

import (
	"reflect"
	"testing"
	"time"

	"github.com/bobinette/deadpool/bandit/components"
)

func TestPlayer(t *testing.T) {
	var tts = map[string]struct {
		arm      int32
		value    float64
		state    components.PlayerState
		expected components.PlayerState
	}{
		"simple state": {
			arm:      0,
			value:    2,
			state:    components.PlayerState{Knowledge: components.Arms{0}, RemainingPlies: 1, Score: 0},
			expected: components.PlayerState{Knowledge: components.Arms{2}, RemainingPlies: 0, Score: 2},
		},
	}

	p := NewPlayer()
	p.Start()
	d := 100 * time.Millisecond
	for name, tt := range tts {
		pushed := false
		select {
		case p.arm.Ch <- tt.arm:
			pushed = true
		case <-time.After(d):
			t.Errorf("%s - Could not send arm, waited for %v", name, d)
		}
		if !pushed {
			continue
		}

		pushed = false
		select {
		case p.stateInput.Ch <- tt.state:
			pushed = true
		case <-time.After(d):
			t.Errorf("%s - Could not send state, waited for %v", name, d)
		}
		if !pushed {
			continue
		}

		pushed = false
		select {
		case p.value.Ch <- tt.value:
			pushed = true
		case <-time.After(d):
			t.Errorf("%s - Could not send value, waited for %v", name, d)
		}
		if !pushed {
			continue
		}

		select {
		case state := <-p.State:
			if !reflect.DeepEqual(state, tt.expected) {
				t.Errorf("%s - Incorrect state: expected %v got %v", name, tt.expected, state)
			}
		case <-time.After(d):
			t.Errorf("%s - Could not receive state, waited for %v", name, d)
		}
	}
}
