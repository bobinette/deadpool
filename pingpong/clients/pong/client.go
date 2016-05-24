package main

import (
	"errors"
	"io"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/pingpong/protos"
)

const name = "Pong"

type Client interface {
	Connect() error
	Disconnect() error
}

type Sound int32

const (
	Ping    Sound = iota
	Pong    Sound = iota
	Silence Sound = iota
)

func (s Sound) String() string {
	switch s {
	case Ping:
		return "Ping"
	case Pong:
		return "Pong"
	case Silence:
		return "Silence"
	}
	return "Unknown"
}

type client struct {
	id    int32
	sound Sound
	ppc   protos.PingPongClient

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
	rep := n.GetConnectReply()
	if rep == nil {
		return errors.New("Error: got nil connect reply from server")
	}

	c.id = rep.Id
	c.sound = Sound(rep.Sound)
	log.Printf("Got id %d, sound %v", c.id, c.sound)

	return c.monitor(stream)
}

func (c *client) Disconnect() error {
	// Close stream
	err := c.stream.CloseSend()
	if err != nil {
		log.Printf("Error closing stream: %v", err)
		return err
	}
	log.Println("Stream closed")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := c.ppc.Leave(ctx, &protos.IdMessage{Id: c.id}); err != nil {
		return err
	}
	return nil
}

func (c *client) monitor(stream protos.PingPong_ConnectClient) error {
	c.stream = stream

	for {
		n, err := c.stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error in monitoring: %v", err)
			return err
		}
		err = c.handle(n)
		if err != nil {
			log.Printf("Error handling notification: %v", err)
			return err
		}
	}
}

func (c *client) handle(n *protos.Notification) error {
	switch body := n.Body.(type) {
	case *protos.Notification_GameStatus:
		if Sound(body.GameStatus.CurrentSound) == c.sound {
			return c.play()
		}
	}
	return nil
}

func (c *client) play() error {
	time.Sleep(3 * time.Second)
	req := protos.PlayRequest{
		Id: c.id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	res, err := c.ppc.Play(ctx, &req)
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
