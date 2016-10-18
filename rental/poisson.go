package rental

import (
	"math"
	"math/rand"
	"time"
)

type Poisson struct {
	Rand *rand.Rand
}

func NewPoisson() *Poisson {
	return &Poisson{Rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// Draw returns a number randomly generated from the Poisson distribution defined by lambda.
// A more perfomant seems to be described in the following paper, but for our use case in
// this game, this one is enough.
//
// http://luc.devroye.org/devroye-poisson.pdf
func (d *Poisson) Draw(lambda int) int {
	L := math.Pow(math.E, -float64(lambda))
	k := 0
	p := 1.0

	for p > L {
		k++
		p *= d.Rand.Float64()
	}
	return k - 1
}

func (d *Poisson) P(lambda, n int) float64 {
	l := math.Pow(float64(lambda), float64(n))
	e := math.Pow(math.E, -float64(lambda))
	f := float64(factorial(n))
	return l / f * e

}

func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}
