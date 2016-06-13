package placers

import (
	"log"
	"math/rand"
	"time"

	"github.com/bobinette/deadpool/battleship/proto"
)

type Random struct{}

func (p *Random) Place() []*proto.Ship {
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
		valid = IsPlacementValid(ships)
		if !valid {
			log.Println("Disposition not valid, retrying")
		}
	}
	return ships
}
