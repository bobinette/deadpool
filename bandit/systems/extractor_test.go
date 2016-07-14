package systems

import (
	"testing"
	"time"

	"github.com/bobinette/deadpool/bandit/components"
)

func TestExtractor(t *testing.T) {
	var tts = map[string]struct {
		bandit   components.Arms
		arm      int32
		expected float64
	}{
		"simple bandit": {
			bandit:   components.Arms{2},
			arm:      0,
			expected: 2,
		},
		"3 arms": {
			bandit:   components.Arms{2, 3, 1},
			arm:      2,
			expected: 1,
		},
	}

	e := NewExtractor()
	e.Start()
	d := 100 * time.Millisecond
	for name, tt := range tts {
		pushed := false
		select {
		case e.arm.Ch <- tt.arm:
			pushed = true
		case <-time.After(d):
			t.Errorf("%s - Could not send arm, waited for %v", name, d)
		}
		if !pushed {
			continue
		}

		pushed = false
		select {
		case e.bandit.Ch <- tt.bandit:
			pushed = true
		case <-time.After(d):
			t.Errorf("%s - Could not send bandit, waited for %v", name, d)
		}
		if !pushed {
			continue
		}

		select {
		case iv := <-e.value.Ch:
			v, ok := iv.(float64)
			if !ok {
				t.Errorf("%s - Incorrect value type: expected %T got %T", name, tt.expected, iv)
			} else if v != tt.expected {
				t.Errorf("%s - Incorrect value: expected %f got %f", name, tt.expected, v)
			}
		case <-time.After(d):
			t.Errorf("%s - Did not receive value, waited for %v", name, d)
		}
	}
}
