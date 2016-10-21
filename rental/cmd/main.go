package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/bobinette/deadpool/rental"
)

func main() {
	// r := flag.Bool("reset", false, "set to true to reset the cache file")
	flag.Parse()

	game := rental.NewGame()
	p := rental.MC{Rand: rand.New(rand.NewSource(time.Now().UnixNano()))}

	p.Play(game)
}
