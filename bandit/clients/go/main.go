package main

import (
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
	c := NewClient(conn)
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
