package rental

type Noop struct{}

func (*Noop) Play(g *Game) {
	for !g.Over() {
		g.Play(Input{Cars: 0})
	}
}
