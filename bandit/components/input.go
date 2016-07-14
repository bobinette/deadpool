package components

type Input struct {
	Arm int32
}

type InputEvent struct {
	Input  Input
	Output chan PlayerState
}
