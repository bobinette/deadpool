package pipeline

import (
	"testing"
	"time"
)

func TestPipeline_Add(t *testing.T) {
	var tts = []struct {
		name string
		unit Unit
		fail bool
	}{
		{
			name: "Add first unit",
			unit: &SimpleUnit{
				inputs:  []*Channel{&Channel{Key: "ch", Ch: nil}},
				outputs: []*Channel{&Channel{Key: "ch", Ch: nil}},
			},
			fail: false,
		},
		{
			name: "Adding another unit with same output",
			unit: &SimpleUnit{
				inputs:  []*Channel{&Channel{Key: "ch", Ch: nil}},
				outputs: []*Channel{&Channel{Key: "ch", Ch: nil}},
			},
			fail: true,
		},
		{
			name: "Unit with no input",
			unit: &SimpleUnit{
				inputs:  []*Channel{},
				outputs: []*Channel{&Channel{Key: "ch2", Ch: nil}},
			},
			fail: false,
		},
		{
			name: "Unit with no output",
			unit: &SimpleUnit{
				inputs:  []*Channel{&Channel{Key: "ch", Ch: nil}},
				outputs: []*Channel{},
			},
			fail: false,
		},
		{
			name: "Unit with multiple outputs",
			unit: &SimpleUnit{
				inputs: []*Channel{&Channel{Key: "ch", Ch: nil}},
				outputs: []*Channel{
					&Channel{Key: "ch3", Ch: nil},
					&Channel{Key: "ch4", Ch: nil},
				},
			},
			fail: false,
		},
		{
			name: "Unit with multiple inputs",
			unit: &SimpleUnit{
				inputs: []*Channel{
					&Channel{Key: "ch", Ch: nil},
					&Channel{Key: "ch2", Ch: nil},
				},
				outputs: []*Channel{&Channel{Key: "ch5", Ch: nil}},
			},
			fail: false,
		},
		{
			name: "Unit with multiple outputs, one already added",
			unit: &SimpleUnit{
				inputs: []*Channel{&Channel{Key: "ch", Ch: nil}},
				outputs: []*Channel{
					&Channel{Key: "ch", Ch: nil},
					&Channel{Key: "ch6", Ch: nil},
				},
			},
			fail: true,
		},
	}

	pipeline := NewPipeline()
	for _, tt := range tts {
		err := pipeline.Add(tt.unit)
		if err != nil && !tt.fail {
			t.Errorf("%s - Should not have failed, got error: %v", tt.name, err)
		} else if err == nil && tt.fail {
			t.Errorf("%s - Should have failed but did not", tt.name)
		}
	}
}

func TestPipeline_Mount(t *testing.T) {
	var tts = map[string]struct {
		inputUnit  *InputUnit
		units      []Unit
		outputUnit *OutputUnit
		input      string
		output     string
	}{
		"No middle unit": {
			inputUnit:  &InputUnit{NewForwardingUnit("", "input")},
			units:      nil,
			outputUnit: &OutputUnit{NewForwardingUnit("input", "")},
			input:      "yolo",
			output:     "yolo",
		},
		"Forwarding middle unit": {
			inputUnit:  &InputUnit{NewForwardingUnit("", "input")},
			units:      []Unit{NewForwardingUnit("input", "output")},
			outputUnit: &OutputUnit{NewForwardingUnit("output", "")},
			input:      "2swag4u",
			output:     "2swag4u",
		},
	}

	d := 100 * time.Millisecond
	for name, tt := range tts {
		pipeline := NewPipeline()
		if err := pipeline.Add(tt.inputUnit); err != nil {
			t.Errorf("%s - Could not add input unit ", name)
		}
		if err := pipeline.Add(tt.outputUnit); err != nil {
			t.Errorf("%s - Could not add output unit ", name)
		}
		for i, b := range tt.units {
			if err := pipeline.Add(b); err != nil {
				t.Errorf("%s - Could not add unit %d ", name, i)
			}
		}

		if err := pipeline.Mount(); err != nil {
			t.Errorf("%s - Error mounting: %v", name, err)
			continue
		}

		pushed := false
		select {
		case tt.inputUnit.inputCh.Ch <- tt.input:
			pushed = true
		case <-time.After(d):
			t.Errorf("%s - Could not push in inputUnit, waited %v", name, d)
		}
		if !pushed {
			continue
		}

		output := ""
		select {
		case ioutput := <-tt.outputUnit.outputCh.Ch:
			var ok bool
			output, ok = ioutput.(string)
			if !ok {
				t.Errorf("%s - Incorrect type, expected string got %T", name, ioutput)
			}
		case <-time.After(d):
			t.Errorf("%s - Did not receive output, waited %v", name, d)
		}

		if output != tt.output {
			t.Errorf("%s - Incorrect output: expected '%s' got '%s'", name, tt.output, output)
		}
	}
}

// ------------------------------------------------------------------------------------------------
// Helpers

// ---- SimpleUnit
type SimpleUnit struct {
	inputs  []*Channel
	outputs []*Channel
}

func (b *SimpleUnit) Inputs() []*Channel {
	return b.inputs
}

func (b *SimpleUnit) Outputs() []*Channel {
	return b.outputs
}

func (b *SimpleUnit) Start() {}

// ---- InputUnit
type InputUnit struct {
	*ForwardingUnit
}

func (b *InputUnit) Inputs() []*Channel {
	return nil
}

// ---- OutputUnit
type OutputUnit struct {
	*ForwardingUnit
}

func (b *OutputUnit) Outputs() []*Channel {
	return nil
}

// ---- ForwardingUnit
type ForwardingUnit struct {
	inputCh  Channel
	outputCh Channel
}

func NewForwardingUnit(input string, output string) *ForwardingUnit {
	return &ForwardingUnit{
		inputCh: Channel{
			Key: input,
			Ch:  make(chan interface{}),
		},
		outputCh: Channel{
			Key: output,
			Ch:  make(chan interface{}),
		},
	}
}

func (b *ForwardingUnit) Inputs() []*Channel {
	return []*Channel{&b.inputCh}
}

func (b *ForwardingUnit) Outputs() []*Channel {
	return []*Channel{&b.outputCh}
}

func (b *ForwardingUnit) Start() {
	go b.run()
}

func (b *ForwardingUnit) run() {
	for {
		i := <-b.inputCh.Ch
		b.outputCh.Ch <- i
	}
}
