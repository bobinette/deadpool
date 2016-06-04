package players

import (
	"github.com/bobinette/deadpool/battleship/proto"
)

type Player interface {
	// Name returns the name of the player. This is used for tracking which
	// AI is playing, so clear names are recommented
	Name() string

	// Disposition returns the position occupied by the ships as determined
	// by the player. It should be an array of 17 integers representing:
	//
	//             Ship | Size
	// -----------------|------
	// Aircraft carrier |   5
	//       Battleship |   4
	//        Submarine |   3
	//        Destroyer |   3
	//      Patrol boat |   2
	// -----------------|------
	//            Total |  17
	Disposition() []*proto.Ship

	// Play asks the player to play given each side of the board. It takes as
	// argument the two sides and returns the position to fire on
	Play() int32
}

func NewPlayer(name string) Player {
	switch name {
	case "random":
		return &RandomPlayer{}
	}
	return nil
}
