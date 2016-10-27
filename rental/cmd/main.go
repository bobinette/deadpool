package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/bobinette/deadpool/rental"
)

type Player interface {
	Play(rental.Game)
}

func main() {
	playerName := flag.String("player", "mc", "player to use")
	flag.Parse()

	game := rental.NewGame()
	var p Player

	if *playerName == "dp" {
		p = &rental.DP{}
		log.Println("Using Dynamic Programming player")
	} else if *playerName == "mc" {
		p = &rental.MC{Rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
		log.Println("Using Monte Carlo player")
	} else {
		log.Fatalln("Unknown player: ", *playerName)
	}

	p.Play(game)
}
