package rental

import (
	"fmt"
	"math"
	"time"
)

type DP struct {
}

// The probability becomes too small to be of any help above 11
const (
	LoopMax        = 11
	ConstantReturn = false
)

const (
	θ = 0.01
	γ = 0.9
)

func (dp *DP) Play(g Game) {
	// Initialize state value
	v := make(StateValue)
	for i := 0; i < g.Params.MaxCars+1; i++ {
		for j := 0; j < g.Params.MaxCars+1; j++ {
			v[State{i, j}] = 0
		}
	}

	stable := false
	p := make(Policy)
	for !stable {
		// Policy evaluation
		Δ := θ + 1
		for i := 0; Δ > θ && i < 100; i++ {
			Δ = 0.0
			vtmp := v
			v = make(StateValue)

			start := time.Now()
			k := 0
			for s, tmp := range vtmp {
				v[s] = expectedReward(s, p[s], vtmp, g.Params)

				δ := math.Abs(tmp-v[s]) / math.Abs(tmp)
				if δ > Δ {
					Δ = δ
				}

				fmt.Print("\r", k, " Δ: ", Δ, "    ", time.Since(start))
				k++
			}
			fmt.Println("\ni:", i, "\t\tΔ:", Δ)
			fmt.Println()
		}
		fmt.Println()
		fmt.Println(v)

		// Policy improvement
		np := make(Policy)
		start := time.Now()
		for s := range v {
			action := 0
			reward := math.Inf(-1)
			for a := -g.Params.MaxMoves; a < g.Params.MaxMoves+1; a++ {
				if s.CarsAt1-a < 0 || s.CarsAt2+a < 0 {
					continue
				}

				r := expectedReward(s, a, v, g.Params)
				fmt.Print("\r", s.CarsAt1, s.CarsAt2, fmt.Sprintf("%2d", a), r, time.Since(start))
				if r > reward {
					action = a
					reward = r
				}
			}
			np[s] = action
		}
		fmt.Println("\r")
		stable = true
		for s, a := range np {
			if a != p[s] {
				stable = false
				fmt.Println("Policy not stable")
				fmt.Println()
				fmt.Println()
				break
			}
		}
		p = np
	}

	fmt.Println(p)
}

func expectedReward(s State, a int, v StateValue, params Parameters) float64 {
	reward := -2.0 * math.Abs(float64(a))

	newState := State{
		CarsAt1: bounded(0, params.MaxCars, s.CarsAt1-a),
		CarsAt2: bounded(0, params.MaxCars, s.CarsAt2+a),
	}
	d := NewPoisson()

	for ncust1 := 0; ncust1 < LoopMax+1; ncust1++ {
		p1 := d.P(params.CustomerAt1, ncust1)
		rental1 := min(ncust1, newState.CarsAt1)

		for ncust2 := 0; ncust2 < LoopMax+1; ncust2++ {
			p2 := d.P(params.CustomerAt2, ncust2)
			rental2 := min(ncust2, newState.CarsAt2)

			rent := 10 * (rental1 + rental2)

			// Constant number of returned car, much quicker
			if ConstantReturn {
				prob := p1 * p2
				finalState := State{
					CarsAt1: min(newState.CarsAt1-rental1+params.ReturnAt1, params.MaxCars),
					CarsAt2: min(newState.CarsAt2-rental2+params.ReturnAt2, params.MaxCars),
				}
				reward += prob * (float64(rent) + γ*v[finalState])
			} else {
				for nret1 := 0; nret1 < LoopMax+1; nret1++ {
					pr1 := d.P(params.ReturnAt1, nret1)

					for nret2 := 0; nret2 < LoopMax+1; nret2++ {
						pr2 := d.P(params.ReturnAt2, nret2)

						prob := p1 * p2 * pr1 * pr2
						finalState := State{
							CarsAt1: min(newState.CarsAt1-rental1+nret1, params.MaxCars),
							CarsAt2: min(newState.CarsAt2-rental2+nret2, params.MaxCars),
						}
						reward += prob * (float64(rent) + γ*v[finalState])
					}
				}
			}

		}
	}

	return reward
}
