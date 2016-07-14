package bandit

import (
	"reflect"
	"testing"
	"time"

	"github.com/bobinette/deadpool/bandit/components"
)

func TestGame_Play(t *testing.T) {
	var tts = []struct {
		arm      int32
		expected components.PlayerState
	}{
		{
			arm: 2,
			expected: components.PlayerState{
				Knowledge:      components.Arms{0, 0, 3, 0},
				RemainingPlies: 3,
				Score:          3,
			},
		},
		{
			arm: 1,
			expected: components.PlayerState{
				Knowledge:      components.Arms{0, 2, 3, 0},
				RemainingPlies: 2,
				Score:          5,
			},
		},
		{
			arm: 2,
			expected: components.PlayerState{
				Knowledge:      components.Arms{0, 2, 3, 0},
				RemainingPlies: 1,
				Score:          8,
			},
		},
		{
			arm: 3,
			expected: components.PlayerState{
				Knowledge:      components.Arms{0, 2, 3, -2},
				RemainingPlies: 0,
				Score:          6,
			},
		},
	}

	gm, err := NewGameManager(4)
	if err != nil {
		t.Fatalf("Could not create game manager: %v", err)
	}
	gm.bandit = components.Arms{1, 2, 3, -2}
	gm.playerState = components.PlayerState{
		Knowledge:      components.Arms{0, 0, 0, 0},
		RemainingPlies: 4,
		Score:          0,
	}
	gm.Start()

	d := 100 * time.Millisecond
	for i, tt := range tts {
		evt := components.InputEvent{
			Input:  components.Input{Arm: tt.arm},
			Output: make(chan components.PlayerState),
		}

		select {
		case gm.InputChannel <- evt:
		case <-time.After(d):
			t.Fatalf("%d - Could not push event input", i)
		}

		select {
		case output := <-evt.Output:
			if !reflect.DeepEqual(output, tt.expected) {
				t.Errorf("%d - Incorrect output: expected %+v got %+v", i, tt.expected, output)
			}
		case <-time.After(d):
			t.Errorf("%d - Did not receive output", i)
		}
	}
}
