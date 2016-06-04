package battleship

import (
	"fmt"
	"log"

	"github.com/bobinette/deadpool/battleship/proto"
)

type Board []proto.Tile

func (b Board) AddShips(ships []*proto.Ship) {
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
	s += "------------\n"
	for i := 0; i < 10; i++ {
		s += "|"
		for j := 0; j < 10; j++ {
			switch b[10*i+j] {
			case proto.Tile_SEA:
				s += " "
			case proto.Tile_UNKNOWN:
				s += "-"
			case proto.Tile_SHIP:
				s += "O"
			case proto.Tile_SUNKEN:
				s += "X"
			}
		}
		s += "|"
		s += "\n"
	}
	s += "------------"
	return s
}

type Game struct {
	ships     map[int32][]*proto.Ship
	shipTiles map[int32][]map[int32]bool
}

func NewGame() *Game {
	return &Game{
		ships:     make(map[int32][]*proto.Ship),
		shipTiles: make(map[int32][]map[int32]bool),
	}
}

func (g *Game) SaveDisposition(c *Client, ships []*proto.Ship) error {
	if len(ships) != 5 {
		return fmt.Errorf("Should have gotten 5 ships, got %d", len(ships))
	}

	shipTiles := make([]map[int32]bool, len(ships))
	size := make(map[int32]int32)
	for i, ship := range ships {
		if err := g.checkShipPosition(ship); err != nil {
			return err
		}
		size[ship.Size] += 1
		tiles := make(map[int32]bool, ship.Size)
		f := 1
		if ship.Vert {
			f = 10
		}
		for j := 0; j < int(ship.Size); j++ {
			tiles[ship.Pos+int32(j*f)] = false
		}
		shipTiles[i] = tiles
	}

	var allowedSizes = []struct {
		size int32
		n    int32
	}{{5, 1}, {4, 1}, {3, 2}, {2, 1}}
	for _, s := range allowedSizes {
		n, ok := size[s.size]
		if !ok {
			return fmt.Errorf("Could not find any ship of size %d", s.size)
		} else if n != s.n {
			return fmt.Errorf("Incorrect number of ships of size %d: expected %d got %d", s.size, s.n, n)
		}
	}
	if len(size) != len(allowedSizes) {
		return fmt.Errorf("Superfluous ships detected")
	}

	g.ships[c.ID] = ships
	g.shipTiles[c.ID] = shipTiles
	board := make(Board, 100)
	board.AddShips(ships)
	log.Println(board.String())
	return nil
}

func (g *Game) RegisterPly(c *Client, p int32) proto.Tile {
	oID := g.opponentID(c)
	if oID == -1 {
		return proto.Tile_UNKNOWN
	}

	for _, ships := range g.shipTiles[oID] {
		_, ok := ships[p]
		if ok {
			ships[p] = true
			for _, s := range ships {
				if !s {
					return proto.Tile_SHIP
				}
			}
			return proto.Tile_SUNKEN
		}
	}

	return proto.Tile_SEA
}

func (g *Game) checkShipPosition(ship *proto.Ship) error {
	if ship.Pos < 0 {
		return fmt.Errorf("Invalid position (< 0) for ship %+v", ship)
	}

	if ship.Vert && (ship.Pos+10*ship.Size > 99) {
		return fmt.Errorf("Invalid position (too far vertically) for ship %+v", ship)
	} else if !ship.Vert && ((ship.Pos%10)+ship.Size > 9) {
		return fmt.Errorf("Invalid position (too far horizontally) for ship %+v", ship)
	}
	return nil
}

func (g *Game) checkOverlapping(ships []proto.Ship) error {
	set := make(map[int32]struct{})
	for _, ship := range ships {
		f := 0
		if ship.Vert {
			f = 10
		}
		for i := 0; i < int(ship.Size); i++ {
			if _, ok := set[ship.Pos+int32(i*f)]; ok {
				return fmt.Errorf("Ships overlap")
			}
			set[ship.Pos+int32(i*f)] = struct{}{}
		}
	}
	return nil
}

func (g *Game) opponentID(c *Client) int32 {
	for id := range g.ships {
		if id != c.ID {
			return id
		}
	}
	return -1
}
