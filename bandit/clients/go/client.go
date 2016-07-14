package main

import (
	"io"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/bandit/clients/go/players"
	"github.com/bobinette/deadpool/bandit/components"
	"github.com/bobinette/deadpool/bandit/proto"
)

type Client struct {
	nArms  int32
	plies  int32
	player players.Player

	bc proto.BanditClient
}

func NewClient(cc *grpc.ClientConn) *Client {
	return &Client{
		player: nil,
		bc:     proto.NewBanditClient(cc),
	}
}

func (c *Client) Connect() error {
	req := &proto.EmptyMessage{}

	stream, err := c.bc.Connect(context.Background(), req)
	if err != nil {
		return err
	}
	log.Println("Connecting, waiting for response...")

	return c.Monitor(stream)
}

func (c *Client) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := c.bc.Disconnect(ctx, &proto.EmptyMessage{})
	if err != nil {
		return err
	}

	log.Println("Disconnection succesful")
	return nil
}

func (c *Client) Monitor(stream proto.Bandit_ConnectClient) error {
	for {
		n, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error in monitoring: %v", err)
			return err
		}

		switch msg := n.Msg.(type) {
		// Save ID
		case *proto.Notification_ConnectReply:
			log.Println("Connected!")
			c.nArms = msg.ConnectReply.NArms
			c.plies = 0
			c.player = players.NewGreedyPlayer(c.nArms)
			log.Printf("Bandit has %d arms", c.nArms)
			return c.Play()
		}
	}
}

func (c *Client) Play() error {
	play := true
	for play {
		log.Printf("======== Ply %d ========", c.plies)

		arm := c.player.Play()
		rep, err := c.bc.Play(context.Background(), &proto.PlayRequest{
			Arm: c.player.Play(),
		})
		if err != nil {
			return err
		}

		log.Printf("Arm %d raised value %f", arm, rep.Knowledge[arm])
		log.Printf("Score: %f", rep.Score)
		c.plies += 1
		play = rep.RemainingPlies > 0

		state := components.PlayerState{
			Knowledge:      rep.Knowledge,
			Score:          rep.Score,
			RemainingPlies: rep.RemainingPlies,
		}
		c.player.Save(state)
	}

	return nil
}
