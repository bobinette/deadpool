package rental

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var (
	firstVisit = true
	newPolicy  = true
)

type MC struct {
	Rand *rand.Rand
}

type ActionValue map[State]map[int]float64
type Returns map[State]map[int][]float64

func (mc MC) Play(g Game) {
	p := make(Policy)
	for i := 0; i < g.Params.MaxCars+1; i++ {
		for j := 0; j < g.Params.MaxCars+1; j++ {
			p[State{i, j}] = mc.Rand.Intn(g.Params.MaxMoves*2+1) - g.Params.MaxMoves
		}
	}

	fmt.Println(p)
	r := make(Returns)

	q := make(ActionValue)
	start := time.Now()
	for i := 0; i < 100000; i++ {
		mc.simulate(g.Params, r, p)

		for s := range r {
			for a := range r[s] {
				if _, ok := q[s]; !ok {
					q[s] = make(map[int]float64)
				}
				if _, ok := q[s][a]; !ok {
					q[s][a] = 0
				}

				q[s][a] = average(r[s][a])
			}

			ps := 0
			f := math.Inf(-1)
			for a := range q[s] {
				if q[s][a] > f {
					ps = a
					f = q[s][a]
				}
			}
			p[s] = ps
		}
		fmt.Print("\r", i, "\t", time.Since(start))
	}

	fmt.Println()
	fmt.Println()
	fmt.Println(p)
}

func (mc MC) simulate(params Parameters, r Returns, p Policy) {
	rewards := make(map[State]map[int]float64)
	poisson := Poisson{Rand: mc.Rand}

	// Play a full episode - every visit MC
	gameOver := false

	// Randomly select first state and first action
	s := State{
		CarsAt1: mc.Rand.Intn(params.MaxCars + 1),
		CarsAt2: mc.Rand.Intn(params.MaxCars + 1),
	}
	a := mc.Rand.Intn(params.MaxMoves*2+1) - params.MaxMoves
	for !gameOver {
		// Play
		reward := -2 * math.Abs(float64(a))
		c1 := poisson.Draw(params.CustomerAt1)
		r1 := poisson.Draw(params.ReturnAt1)

		c2 := poisson.Draw(params.CustomerAt2)
		r2 := poisson.Draw(params.ReturnAt2)
		if c1 > s.CarsAt1 || c2 > s.CarsAt2 {
			gameOver = true
		}

		reward += 10 * float64(min(c1, s.CarsAt1)+min(c2, s.CarsAt2))

		s.CarsAt1 = bounded(0, params.MaxCars, s.CarsAt1-a-c1+r1)
		s.CarsAt2 = bounded(0, params.MaxCars, s.CarsAt2+a-c2+r2)

		if _, ok := rewards[s]; !ok {
			rewards[s] = make(map[int]float64)
		}
		if _, ok := rewards[s][a]; !ok {
			rewards[s][a] = 0
		} else if firstVisit {
			reward = 0
		}

		// Σ γ^k * Rk = R0 + γ(R1 + γ(R2 + ...))
		for os := range rewards {
			for oa := range rewards[s] {
				rewards[os][oa] = γ*rewards[os][oa] + reward
			}
		}

		// Select next action following policy
		a = p[s]
	}

	// Add episode returns to the list of all the returns (per (state, action))
	for s := range rewards {
		for a := range rewards[s] {
			if _, ok := r[s]; !ok {
				r[s] = make(map[int][]float64)
			}
			if _, ok := r[s][a]; !ok {
				r[s][a] = make([]float64, 0)
			}

			r[s][a] = append(r[s][a], rewards[s][a])
		}
	}
}

func average(ns []float64) float64 {
	if len(ns) == 0 {
		return 0
	}

	var s float64 = 0
	for _, n := range ns {
		s += n
	}
	return s / float64(len(ns))
}
