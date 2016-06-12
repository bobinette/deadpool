package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/battleship"
	"github.com/bobinette/deadpool/pingpong"
)

const (
	port = ":17000"
)

func main() {
	// The servers might register their own flags
	flag.Parse()

	var s *grpc.Server

	game := os.Args[len(os.Args)-1] // Game name is the last command line arg
	switch game {
	case "battleship":
		s = battleship.NewServer()
	case "pingpong":
		s = pingpong.NewServer()
	default:
		log.Printf("Nice try! But game %s is not available", game)
		return
	}

	log.Println("Starting server...")
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

	log.Printf("Listening on %s...", lis.Addr())
	s.Serve(lis)
}
