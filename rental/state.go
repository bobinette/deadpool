package rental

import (
	"fmt"
)

type State struct {
	// CarsAtX represents the number of cars available for renting at location X
	CarsAt1 int
	CarsAt2 int
}

type Policy map[State]int

func (p Policy) Size() (int, int) {
	w, h := 0, 0
	for s := range p {
		if s.CarsAt1 > h {
			h = s.CarsAt1
		}

		if s.CarsAt2 > w {
			w = s.CarsAt2
		}
	}
	return w, h
}

func (p Policy) At(x, y int) string {
	return fmt.Sprintf(" %2d ", p[State{x, y}])
}

func (p Policy) String() string {
	return Print2D(p)
}

type StateValue map[State]float64

func (v StateValue) Size() (int, int) {
	w, h := 0, 0
	for s := range v {
		if s.CarsAt1 > h {
			h = s.CarsAt1
		}

		if s.CarsAt2 > w {
			w = s.CarsAt2
		}
	}
	return w, h
}

func (v StateValue) At(x, y int) string {
	return fmt.Sprintf(" %7.3f ", v[State{x, y}])
}

func (v StateValue) String() string {
	return Print2D(v)
}
