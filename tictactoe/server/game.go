package server

import (
	"github.com/bobinette/deadpool/tictactoe/protos"
)

type Shape int32

const (
	E Shape = iota // Empty
	X Shape = iota
	O Shape = iota
)

func (s Shape) String() string {
	switch s {
	case X:
		return "X"
	case O:
		return "O"
	}
	return "-"
}

type Game struct {
	board []Shape
}

func NewGame() *Game {
	return &Game{
		board: make([]Shape, 9),
	}
}

// Board returns the current state of the board as a simple array of Shapes. The
// indexes are as follows:
// 0   1   2
// 3   4   5
// 6   7   8
func (g Game) Board() []Shape {
	return g.board
}

func (g *Game) Play(s Shape, pos int) protos.PlayReply_Status {
	if pos < 0 || pos > len(g.board) {
		return protos.PlayReply_INVALID_POSITION
	}

	if g.board[pos] != E {
		return protos.PlayReply_POSITION_TAKEN
	}

	g.board[pos] = s
	return protos.PlayReply_VALID
}

func (g Game) Winner() Shape {
	// Winning combinations
	wc := []struct {
		a int
		b int
		c int
	}{
		// rows
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		// columns
		{0, 3, 6},
		{1, 4, 7},
		{2, 5, 8},
		// diagonals
		{0, 4, 8},
		{2, 4, 6},
	}

	for _, c := range wc {
		if g.board[c.a] == E {
			continue
		}

		if g.board[c.a] == g.board[c.b] && g.board[c.a] == g.board[c.c] {
			return g.board[c.a]
		}
	}
	return E
}

func (g Game) Draw() bool {
	return g.Winner() == E && g.IsFull()
}

func (g Game) IsFinished() bool {
	return g.Winner() != E || g.IsFull()
}

func (g Game) IsFull() bool {
	for _, s := range g.board {
		if s == E {
			return false
		}
	}
	return true
}
