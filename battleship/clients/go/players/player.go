package players

import (
	"github.com/bobinette/deadpool/battleship/proto"
)

type Player interface {
	// Play asks the player to play given each side of the board. It takes as
	// argument the two sides and returns the position to fire on
	Play() int32

	// SaveResult allows the client to remember the status of the tiles after
	// a ply.
	SaveResult(int32, proto.Tile)
}

func NewPlayer(name string) Player {
	switch name {
	case "random":
		return &Random{
			board: make([]proto.Tile, 100),
			n:     0,
		}
	case "distance_based":
		return NewDistanceBased(manhattan)
	}
	return nil
}
