package battleship

import (
	"reflect"
	"testing"

	"github.com/bobinette/deadpool/battleship/proto"
)

func TestGame_RegisterPly(t *testing.T) {
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

	var pID int32 = 2
	game := Game{
		ships: map[int32][]*proto.Ship{1: nil, 2: nil},
	}
	for name, tt := range tts {
		game.shipTiles = map[int32][]map[int32]bool{
			1: tt.shipTiles,
		}

		if r := game.RegisterPly(pID, tt.p); r != tt.expected {
			t.Errorf("%s - Incorrect result: expected %d got %d", name, tt.expected, r)
		}
	}
}

func TestGame_Winner(t *testing.T) {
	var tts = map[string]struct {
		shipTiles map[int32][]map[int32]bool
		w         int32
	}{
		"Winner is 1": {
			map[int32][]map[int32]bool{
				1: []map[int32]bool{
					{17: true, 18: true, 19: true},
					{21: true, 22: true, 23: true, 24: false},
				},
				2: []map[int32]bool{
					{33: true, 43: true, 53: true},
					{81: true, 82: true},
				},
			},
			1,
		},
		"No winner": {
			map[int32][]map[int32]bool{
				1: []map[int32]bool{
					{17: true, 18: true, 19: true},
					{21: true, 22: true, 23: true, 24: false},
				},
				2: []map[int32]bool{
					{33: true, 43: true, 53: true},
					{81: false, 82: true},
				},
			},
			-1,
		},
	}

	game := Game{
		ships: map[int32][]*proto.Ship{1: nil, 2: nil},
	}
	for name, tt := range tts {
		game.shipTiles = tt.shipTiles
		if w := game.Winner(); w != tt.w {
			t.Errorf("%s - Incorrect winner: expected %d got %d", name, tt.w, w)
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

func TestGame_checkOverlapping(t *testing.T) {
	var tts = map[string]struct {
		ships []*proto.Ship
		valid bool
	}{
		"One square": {
			[]*proto.Ship{
				&proto.Ship{Pos: 5, Size: 1, Vert: false},
				&proto.Ship{Pos: 5, Size: 1, Vert: false},
			},
			false,
		},
		"Full set ok": {
			[]*proto.Ship{
				&proto.Ship{Pos: 3, Size: 5, Vert: false},
				&proto.Ship{Pos: 25, Size: 4, Vert: true},
				&proto.Ship{Pos: 62, Size: 3, Vert: true},
				&proto.Ship{Pos: 11, Size: 3, Vert: false},
				&proto.Ship{Pos: 47, Size: 2, Vert: false},
			},
			true,
		},
	}

	game := Game{}
	for name, tt := range tts {
		err := game.checkOverlapping(tt.ships)
		board := Board(make([]proto.Tile, 100))
		board.AddShips(tt.ships)
		if tt.valid && err != nil {
			t.Errorf("%s - Ships not overlapping but got error %v", name, err)
			t.Log(board.String())
		} else if !tt.valid && err == nil {
			t.Errorf("%s - Ships overlapping but got no error ", name)
			t.Log(board.String())
		}
	}
}

func TestGame_shipTiles(t *testing.T) {
	var tts = map[string]struct {
		ship  proto.Ship
		tiles []int32
	}{
		"Single tile": {
			proto.Ship{Pos: 6, Size: 1, Vert: false},
			[]int32{6},
		},
		"4 tiles vertical": {
			proto.Ship{Pos: 52, Size: 4, Vert: true},
			[]int32{52, 62, 72, 82},
		},
	}

	game := Game{}
	for name, tt := range tts {
		tiles := game.tilesForShip(tt.ship)
		if !reflect.DeepEqual(tiles, tt.tiles) {
			t.Errorf("%s - Invalid tiles from ship: expected %v got %v", name, tt.tiles, tiles)
		}
	}
}
