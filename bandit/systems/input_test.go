package systems

import (
	"reflect"
	"testing"
	"time"

	"github.com/bobinette/deadpool/bandit/components"
)

func TestInput(t *testing.T) {
	var tts = map[string]struct {
		input  *Input
		bandit components.Arms
		state  components.PlayerState
		arm    int32
	}{
		"basic inputs": {
			input:  NewInput(),
			bandit: components.Arms{2},
			state: components.PlayerState{
				Knowledge:      components.Arms{0},
				RemainingPlies: 1,
				Score:          0,
			},
			arm: 0,
		},
	}

	d := 100 * time.Millisecond
	for name, tt := range tts {
		tt.input.Start()

		select {
		case tt.input.Bandit <- tt.bandit:
			select {
			case ibandit := <-tt.input.bandit.Ch:
				bandit, ok := ibandit.(components.Arms)
				if !ok {
					t.Errorf("%s - Incorrect bandit type: expected %T got %T", name, tt.bandit, bandit)
				} else if !reflect.DeepEqual(tt.bandit, bandit) {
					t.Errorf("%s - Incorrect bandit: expected %v got %v", name, tt.bandit, bandit)
				}
			case <-time.After(d):
				t.Errorf("%s - Did not receive bandit, waited for %v", name, d)

			}
		case <-time.After(d):
			t.Errorf("%s - Could not send bandit, waited for %v", name, d)
		}

		select {
		case tt.input.PlayerState <- tt.state:
			select {
			case istate := <-tt.input.playerState.Ch:
				state, ok := istate.(components.PlayerState)
				if !ok {
					t.Errorf("%s - Incorrect state type: expected %T got %T", name, tt.state, state)
				} else if !reflect.DeepEqual(tt.state, state) {
					t.Errorf("%s - Incorrect state: expected %v got %v", name, tt.state, state)
				}
			case <-time.After(d):
				t.Errorf("%s - Did not receive state, waited for %v", name, d)

			}
		case <-time.After(d):
			t.Errorf("%s - Could not send state, waited for %v", name, d)
		}

		select {
		case tt.input.Arm <- tt.arm:
			select {
			case iarm := <-tt.input.arm.Ch:
				arm, ok := iarm.(int32)
				if !ok {
					t.Errorf("%s - Incorrect arm type: expected %T got %T", name, tt.arm, arm)
				} else if !reflect.DeepEqual(tt.arm, arm) {
					t.Errorf("%s - Incorrect arm: expected %v got %v", name, tt.arm, arm)
				}
			case <-time.After(d):
				t.Errorf("%s - Did not receive arm, waited for %v", name, d)

			}
		case <-time.After(d):
			t.Errorf("%s - Could not send arm, waited for %v", name, d)
		}
	}
}
