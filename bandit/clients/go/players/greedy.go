package players

type GreedyPlayer struct {
	Values []float64
}

func NewGreedyPlayer(nArms int32) *GreedyPlayer {
	return &GreedyPlayer{
		Values: make([]float64, nArms),
	}
}

func (p *GreedyPlayer) Play() int32 {
	max := p.Values[0]
	var iMax int32 = 0
	for i, v := range p.Values {
		if v > max {
			max = v
			iMax = int32(i)
		}
	}
	return iMax
}

func (p *GreedyPlayer) SaveValue(i int32, v float64) {
	p.Values[i] = v
}
