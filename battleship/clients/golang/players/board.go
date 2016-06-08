package players

import (
	"fmt"

	"github.com/bobinette/deadpool/battleship/proto"
)

type Board []proto.Tile

func (b Board) PlaceShips(ships []*proto.Ship) {
	for i := 0; i < len(b); i++ {
		b[i] = proto.Tile_SEA
	}

	for _, ship := range ships {
		f := 1
		if ship.Vert {
			f = 10
		}
		for i := 0; i < int(ship.Size); i++ {
			b[ship.Pos+int32(i*f)] = proto.Tile_SHIP
		}
	}
}

func (b Board) String() string {
	s := "\n"
	s += "--0123456789-\n"
	for i := 0; i < 10; i++ {
		s += fmt.Sprintf("%d|", i)
		for j := 0; j < 10; j++ {
			switch b[10*i+j] {
			case proto.Tile_SEA:
				s += "~"
			case proto.Tile_UNKNOWN:
				s += "."
			case proto.Tile_SHIP:
				s += "O"
			case proto.Tile_SUNK:
				s += "X"
			}
		}
		s += "|"
		s += "\n"
	}
	s += "-------------"
	return s
}

// ----------------------------------------------------------------------------
// Check disposition
func IsDispositionValid(ships []*proto.Ship) bool {
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
