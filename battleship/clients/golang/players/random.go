package players

import (
	"log"
	"math/rand"
	"time"

	"github.com/bobinette/deadpool/battleship/proto"
)

type Random struct {
	board Board

	n int
}

func (p *Random) Name() string {
	return "Random"
}

func (p *Random) Play() int32 {
	tiles := make([]int, 0, 100)
	for i, t := range p.board {
		if t == proto.Tile_UNKNOWN {
			tiles = append(tiles, i)
		}
	}

	p.n += 1
	log.Printf("Play n. %d", p.n)

	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)
	return int32(tiles[gen.Intn(len(tiles))])
}

func (p *Random) SaveResult(pos int32, tile proto.Tile) {
	p.board[pos] = tile
	log.Println(p.board.String())
}
