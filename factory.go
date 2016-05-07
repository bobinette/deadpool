package main

import (
	"fmt"

	"google.golang.org/grpc"

	helloworld "github.com/bobinette/deadpool/helloworld/server"
)

func NewServer(game string) (*grpc.Server, error) {
	if game == "helloworld" {
		return helloworld.NewServer(), nil
	}

	return nil, fmt.Errorf("No server available for game %s", game)
}

func AvailableGames() []string {
	return []string{
		"helloworld",
	}
}
