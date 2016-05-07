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
	defer conn.Close()

	c := NewClient(conn)

	// Capture interrupt
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		err := c.Connect()
		if err != nil {
			log.Println(err)
		}
		close(quit)
	}()

	<-quit
	if err := c.Leave(); err != nil {
		log.Printf("Error when leaving: %v", err)
	}
	return
}
