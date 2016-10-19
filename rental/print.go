package rental

import (
	"bytes"
	"fmt"
	"strings"
)

type Axis int

type Printable2D interface {
	Size() (w int, h int)
	At(x, y int) string
}

func Print2D(p Printable2D) string {
	w, h := p.Size()

	buf := bytes.Buffer{}
	wLabel := make(map[int]int)
	for y := h; y >= 0; y-- {
		buf.WriteString(fmt.Sprintf("%2d | ", y))
		for x := 0; x <= w; x++ {
			s := p.At(x, y)
			buf.WriteString(fmt.Sprintf(" %s ", s))

			if len(s) > wLabel[x] {
				wLabel[x] = len(s)
			}
		}
		buf.WriteRune('\n')
	}

	buf.WriteString("----")
	for x := 0; x <= w; x++ {
		buf.WriteString(strings.Repeat("-", wLabel[x]+2))
	}
	buf.WriteRune('\n')
	buf.WriteString("2/1|")
	for x := 0; x <= w; x++ {
		f := fmt.Sprintf(" %%%dd ", wLabel[x])
		buf.WriteString(fmt.Sprintf(f, x))
	}
	return buf.String()
}
