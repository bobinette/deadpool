package players

import (
	"log"
	"math/rand"
	"time"

	"github.com/bobinette/deadpool/battleship/proto"
)

type RandomPlayer struct {
	board Board
}

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
	tiles := make([]int, 0, 100)
	for i, t := range p.board {
		if t == proto.Tile_UNKNOWN {
			tiles = append(tiles, i)
		}
	}

	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)

	return int32(tiles[gen.Intn(len(tiles))])
}

func (p RandomPlayer) SaveResult(pos int32, tile proto.Tile) {
	p.board[pos] = tile
	log.Println(p.board.String())
}
