package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "localhost:17000"
)

func main() {
	player := flag.String("player", "random", "player that will be used for the game")
	disposer := flag.String("disposer", "random", "disposer that will be used for the game")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBackoffMaxDelay(2*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing the connection: %v", err)
		}
	}()

	// Capture interrupt
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	log.Println("Dialing...")
	c := NewClient(conn, *player, *disposer)
	go func() {
		err := c.Connect()
		if err != nil {
			log.Printf("Error: %s", err)
		}
		close(quit)
	}()
	defer c.Disconnect()

	<-quit
	log.Println("Done")
}
