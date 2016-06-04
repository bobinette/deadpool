package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	helloworld "github.com/bobinette/deadpool/helloworld/server"
)

func Up(game string) error {
	s, err := NewServer(game)
	if err != nil {
		return err
	}

	log.Printf("Starting server for game %s...", game)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Handle CTRL-C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		s.Stop()
	}()

	log.Println("Listening...")
	return s.Serve(lis)
}

func NewServer(game string) (*grpc.Server, error) {
	switch(game) {
	case "helloworld":
		return helloworld.NewServer(), nil
	}

	return nil, fmt.Errorf("No server available for game %s", game)
}

func AvailableGames() []string {
	return []string{
		"connect4",
		"helloworld",
		"tictactoe",
	}
}
