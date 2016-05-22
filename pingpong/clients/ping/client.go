package main

import (
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/pingpong/protos"
)

const name = "Ping"

type Client interface {
	Connect() error
	Disconnect() error
}

type client struct {
	id  int32
	ppc protos.PingPongClient

	stream protos.PingPong_ConnectClient
}

func NewClient(cc *grpc.ClientConn) Client {
	return &client{
		ppc: protos.NewPingPongClient(cc),
	}
}

func (c *client) Connect() error {
	stream, err := c.ppc.Connect(context.Background(), &protos.ConnectRequest{Name: name})
	if err != nil {
		return err
	}
	log.Println("Successfully contacted server")

	// Wait for id
	n, err := stream.Recv()
	if err != nil {
		return err
	}
	c.id = n.GetIdMessage().Id
	log.Printf("Got id %d", c.id)

	go c.monitor(stream)

	go func() {
		src := rand.NewSource(time.Now().UnixNano())
		gen := rand.New(src)
		for {
			time.Sleep(time.Duration(gen.Intn(5)) * time.Second)
			err := c.play()
			if err != nil {
				log.Printf("Error while requesting for game status: %v", err)
				return
			}
		}
	}()

	return nil
}

func (c *client) Disconnect() error {
	// Close stream
	err := c.stream.CloseSend()
	if err != nil {
		log.Printf("Error closing stream: %v", err)
		return err
	}
	log.Println("Stream closed")
	return nil
}

func (c *client) monitor(stream protos.PingPong_ConnectClient) {
	c.stream = stream

	// Randomly request for game status
	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)
	for {
		time.Sleep(time.Duration(gen.Intn(20)) * time.Second)
		req := protos.IdMessage{Id: c.id}
		s, err := c.ppc.GetGameStatus(context.Background(), &req)
		if err != nil {
			log.Printf("Error while requesting for game status: %v", err)
			return
		}
		log.Println(s.Pingpong)
	}
}

func (c *client) play() error {
	req := protos.PlayRequest{
		Id: c.id,
	}
	res, err := c.ppc.Play(context.Background(), &req)
	if err != nil {
		return err
	}

	if !res.Accepted {
		log.Println("Play not accepted")
	} else {
		log.Println("Play was accepted")
	}
	return nil
}
