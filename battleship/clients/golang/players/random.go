package players

import (
	"github.com/bobinette/deadpool/battleship/proto"
)

type RandomPlayer struct{}

func (p RandomPlayer) Name() string {
	return "Random"
}

func (p RandomPlayer) Disposition() []*proto.Ship {
	sizes := []int32{5, 4, 3, 3, 2}
	ships := make([]*proto.Ship, len(sizes))
	for i, s := range sizes {
		ships[i] = &proto.Ship{
			Pos:  10 * int32(i),
			Vert: false,
			Size: s,
		}
	}
	return ships
}

func (p RandomPlayer) Play() int32 {
	return 0
}
