package players

import (
	"log"
	"math/rand"

	"github.com/bobinette/deadpool/battleship/proto"
)

type Pos struct {
	x int32
	y int32
}

func FromLinear(p int32) Pos {
	return Pos{
		x: p % 10,
		y: p / 10,
	}
}

func (p Pos) ToLinear() int32 {
	return 10*p.y + p.x
}

type Distance func(p1, p2 Pos) float64

func manhattan(p1, p2 Pos) float64 {
	dx := p2.x - p1.x
	if dx < 0 {
		dx = -dx
	}

	dy := p2.y - p1.y
	if dy < 0 {
		dy = -dy
	}

	return float64(dx + dy)
}

type DistanceBased struct {
	board  Board
	fired  []Pos
	dist   Distance
	scores map[int32]float64

	n int
}

func NewDistanceBased(dist Distance) Player {
	scores := make(map[int32]float64)
	var i int32
	for i = 0; i < 100; i++ {
		pos := FromLinear(i)
		d := dist(pos, Pos{pos.x, -1})
		dd := dist(pos, Pos{-1, pos.y})
		if dd < d {
			d = dd
		}

		dd = dist(pos, Pos{pos.x, 10})
		if dd < d {
			d = dd
		}

		dd = dist(pos, Pos{10, pos.y})
		if dd < d {
			d = dd
		}

		scores[i] = d
	}
	return &DistanceBased{
		board:  make([]proto.Tile, 100),
		fired:  make([]Pos, 0, 100),
		dist:   dist,
		scores: scores,

		n: 0,
	}
}

func (p *DistanceBased) Play() int32 {
	pos := -1
	var max float64 = 0
	for i, t := range p.board {
		if t != proto.Tile_UNKNOWN {
			continue
		}

		s := p.scores[int32(i)]
		if s > max || pos == -1 {
			pos = i
			max = s
		}
	}

	p.n += 1
	log.Printf("Play n. %d", p.n)
	return int32(pos)
}

func (p *DistanceBased) SaveResult(pos int32, tile proto.Tile) {
	p.board[pos] = tile
	p.fired = append(p.fired, FromLinear(pos))

	var max float64 = 0
	for t, s := range p.scores {
		p1 := FromLinear(t)
		p2 := FromLinear(pos)
		d := p.dist(p1, p2) + rand.NormFloat64()
		if d < s {
			p.scores[t] = d
		}

		if p.scores[t] > max {
			max = p.scores[t]
		}
	}

	for i, t := range p.board {
		if t != proto.Tile_SHIP {
			continue
		}

		for _, a := range p.adjacentTiles(int32(i)) {
			if p.board[a] == proto.Tile_UNKNOWN {
				p.scores[a] = max
			}
		}
	}

	log.Println(p.board.String())
}

func (p *DistanceBased) adjacentTiles(i int32) []int32 {
	adj := make([]int32, 0, 4)

	if i >= 0 {
		adj = append(adj, i-1)
	}
	if i%10 < 9 {
		adj = append(adj, i+1)
	}
	if i/10 > 0 {
		adj = append(adj, i-10)
	}
	if i/10 < 9 {
		adj = append(adj, i+10)
	}

	return adj
}
