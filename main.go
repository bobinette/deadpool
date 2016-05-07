package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/codegangsta/cli"
)

const (
	port = ":17000"
)

func startServer(game string) error {
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

func main() {
	app := cli.NewApp()
	app.Name = "deadpool"
	app.Usage = "let the AIs fight"

	app.Commands = []cli.Command{
		{
			Name:    "up",
			Aliases: []string{"u"},
			Usage:   "start a deadpool server",
			Action: func(c *cli.Context) error {
				err := startServer(c.Args().First())
				if err != nil {
					log.Println(err)
				}
				return err
			},
		},
	}

	app.Run(os.Args)
}
