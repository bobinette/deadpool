package battleship

import (
	"testing"

	"github.com/bobinette/deadpool/battleship/proto"
)

func TestGame_RegisterPly(t *testing.T) {
	c := Client{ID: 2}
	game := Game{
		ships: map[int32][]*proto.Ship{1: nil, 2: nil},
	}
	var tts = map[string]struct {
		shipTiles []map[int32]bool
		p         int32
		expected  proto.Tile
	}{
		"One tile ship sunken": {
			[]map[int32]bool{
				map[int32]bool{
					17: false,
				},
			},
			17,
			proto.Tile_SUNKEN,
		},
		"Two tile ship hit": {
			[]map[int32]bool{
				map[int32]bool{
					17: false,
					18: false,
				},
			},
			17,
			proto.Tile_SHIP,
		},
		"Sea": {
			[]map[int32]bool{
				map[int32]bool{
					17: false,
					18: false,
				},
			},
			5,
			proto.Tile_SEA,
		},
		"Three tile ship sunken": {
			[]map[int32]bool{
				map[int32]bool{
					7:  true,
					17: false,
					27: true,
				},
			},
			17,
			proto.Tile_SUNKEN,
		},
	}

	for name, tt := range tts {
		game.shipTiles = map[int32][]map[int32]bool{
			1: tt.shipTiles,
		}

		if r := game.RegisterPly(&c, tt.p); r != tt.expected {
			t.Errorf("%s - Incorrect result: expected %d got %d", name, tt.expected, r)
		}
	}
}

func TestGame_checkShipPosition(t *testing.T) {
	game := Game{}

	var tts = map[string]struct {
		ship  proto.Ship
		valid bool
	}{
		"Invalid < 0": {
			proto.Ship{Size: 1, Pos: -3, Vert: false},
			false,
		},
		"Invalid too far horizontally": {
			proto.Ship{Size: 5, Pos: 5, Vert: false},
			false,
		},
		"Invalid too far horizontally - row 6": {
			proto.Ship{Size: 4, Pos: 57, Vert: false},
			false,
		},
		"Invalid too far vertically": {
			proto.Ship{Size: 5, Pos: 60, Vert: true},
			false,
		},
		"Invalid too far vertically - col 6": {
			proto.Ship{Size: 5, Pos: 56, Vert: true},
			false,
		},
		"Valid - row 6": {
			proto.Ship{Size: 4, Pos: 53, Vert: false},
			true,
		},
		"Valid - col 6": {
			proto.Ship{Size: 5, Pos: 46, Vert: true},
			true,
		},
	}

	for name, tt := range tts {
		err := game.checkShipPosition(&tt.ship)
		if !tt.valid && err == nil {
			t.Errorf("%s - Ship %+v not valid but check raised no error", name, tt.ship)
		} else if tt.valid && err != nil {
			t.Errorf("%s - Ship %+v valid but got error %v", name, tt.ship, err)
		}
	}
}
