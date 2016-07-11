package bandit

type Bandit struct {
	Values []float64
}

type Player struct {
	History        []int32
	Score          float64
	RemainingPlies int32
	Knowledge      []float64
}

type Input struct {
	Arm int32
}
