package placers

import (
	"github.com/bobinette/deadpool/battleship/proto"
)

type Placer interface {
	// Placement returns the position occupied by the ships as determined
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
	Place() []*proto.Ship
}

func NewPlacer(name string) Placer {
	switch name {
	case "random":
		return &Random{}
	}
	return nil
}

// ----------------------------------------------------------------------------
// Check Placement
func IsPlacementValid(ships []*proto.Ship) bool {
	set := make(map[int32]int32)
	for _, ship := range ships {
		f := 1
		if ship.Vert {
			f = 10
		}
		for i := 0; i < int(ship.Size); i++ {
			if _, ok := set[ship.Pos+int32(i*f)]; ok {
				return false
			}
			set[ship.Pos+int32(i*f)] = ship.Size
		}
	}
	return true
}
