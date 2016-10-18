package main

import (
	"flag"
	"github.com/bobinette/deadpool/rental"
)

func main() {
	// r := flag.Bool("reset", false, "set to true to reset the cache file")
	flag.Parse()

	game := rental.NewGame()
	p := rental.DP{}

	p.Play(game)
}
