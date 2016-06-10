package main

import (
	"errors"
	"io"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/pingpong/clients/golang/players"
	"github.com/bobinette/deadpool/pingpong/proto"
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
	id     int32
	sound  Sound
	player players.Player

	ppc proto.PingPongClient
}

func NewClient(cc *grpc.ClientConn, player string) Client {
	return &client{
		player: players.NewPlayer(player),
		ppc:    proto.NewPingPongClient(cc),
	}
}

func (c *client) Connect() error {
	stream, err := c.ppc.Connect(context.Background(), &proto.ConnectRequest{Name: name})
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

	return c.Monitor(stream)
}

func (c *client) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := c.ppc.Disconnect(ctx, &proto.IdMessage{Id: c.id}); err != nil {
		return err
	}

	log.Println("Disconnection succesful")
	return nil
}

func (c *client) Monitor(stream proto.PingPong_ConnectClient) error {
	for {
		n, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error in monitoring: %v", err)
			return err
		}
		err = c.Handle(n)
		if err != nil {
			log.Printf("Error handling notification: %v", err)
			return err
		}
	}
}

func (c *client) Handle(n *proto.Notification) error {
	switch body := n.Body.(type) {
	case *proto.Notification_GameStatus:
		if Sound(body.GameStatus.CurrentSound) == c.sound {
			return c.Play()
		}
	}
	return nil
}

func (c *client) Play() error {
	c.player.Play()
	req := proto.PlayRequest{
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
