package pipeline

type Unit interface {
	Inputs() []*Channel
	Outputs() []*Channel

	Start()
}
