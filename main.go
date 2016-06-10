package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bobinette/deadpool/pingpong"
)

const (
	port = ":17000"
)

func main() {
	s := pingpong.NewServer()

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
