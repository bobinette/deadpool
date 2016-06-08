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

func (p *Random) Disposition() []*proto.Ship {
	sizes := []int{5, 4, 3, 3, 2}
	ships := make([]*proto.Ship, len(sizes))

	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)

	valid := false
	for !valid {
		for i, s := range sizes {
			major := gen.Intn(10 - s)
			minor := gen.Intn(10)
			v := gen.Intn(2) == 0

			p := 0
			if v {
				p = 10*major + minor
			} else {
				p = major + 10*minor
			}
			ships[i] = &proto.Ship{
				Pos:  int32(p),
				Vert: v,
				Size: int32(s),
			}
		}
		valid = IsDispositionValid(ships)
		if !valid {
			log.Println("Disposition not valid, retrying")
		}
	}
	board := Board(make([]proto.Tile, 100))
	board.PlaceShips(ships)
	log.Printf("Settled for:%s", board.String())
	return ships
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
