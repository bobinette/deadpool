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
	sizes := []int{5, 4, 3, 3, 2}
	ships := make([]*proto.Ship, len(sizes))

	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)

	overlap := true
	for overlap {
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
		overlap = p.checkOverlapping(ships)
		if overlap {
			log.Println("Disposition overlapped, retrying")
		}
	}
	board := Board(make([]proto.Tile, 100))
	board.PlaceShips(ships)
	log.Printf("Settled for:%s", board.String())
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

func (p RandomPlayer) checkOverlapping(ships []*proto.Ship) bool {
	set := make(map[int32]int32)
	for _, ship := range ships {
		f := 1
		if ship.Vert {
			f = 10
		}
		for i := 0; i < int(ship.Size); i++ {
			if _, ok := set[ship.Pos+int32(i*f)]; ok {
				return true
			}
			set[ship.Pos+int32(i*f)] = ship.Size
		}
	}
	return false
}
