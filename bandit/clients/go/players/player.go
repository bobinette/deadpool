package players

type Player interface {
	Play() int32
	SaveValue(int32, float64)
}
