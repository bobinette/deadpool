package players

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/bobinette/deadpool/bandit/components"
)

var (
	epsilon float64 = 0.1
)

func init() {
	flag.Float64Var(&epsilon, "epsilon", epsilon, "decrease rate")
}

type GreedyPlayer struct {
	values  []float64
	gen     *rand.Rand
	epsilon float64
}

func NewGreedyPlayer(nArms int32) *GreedyPlayer {
	src := rand.NewSource(time.Now().UnixNano())
	return &GreedyPlayer{
		values:  make([]float64, nArms),
		gen:     rand.New(src),
		epsilon: epsilon,
	}
}

func (p *GreedyPlayer) Play() int32 {
	max := p.values[0]
	var iMax int32 = 0
	for i, v := range p.values {
		if v > max {
			max = v
			iMax = int32(i)
		}
	}
	log.Printf("max is %f at arm %d", max, iMax)

	e := p.gen.Float64()
	log.Printf("epsilon %f - gen %f", p.epsilon, e)
	if e < p.epsilon {
		iMax = p.gen.Int31n(int32(len(p.values)))
		log.Printf("selecting %d (value: %f) instead of max", iMax, p.values[iMax])
	}

	return iMax
}

func (p *GreedyPlayer) Save(s components.PlayerState) {
	p.values = s.Knowledge
}
