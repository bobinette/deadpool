package components

type Arms []float64

type PlayerState struct {
	Knowledge      Arms
	Score          float64
	RemainingPlies int32
}
