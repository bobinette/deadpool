package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/battleship/clients/golang/players"
	"github.com/bobinette/deadpool/battleship/proto"
)

type Client struct {
	ID     int32
	Player players.Player

	bc proto.BattleshipClient
}

func NewClient(cc *grpc.ClientConn, player string) *Client {
	return &Client{
		bc:     proto.NewBattleshipClient(cc),
		Player: players.NewPlayer(player),
	}
}

func (c *Client) Connect() error {
	if c.Player == nil {
		return fmt.Errorf("No player selected, cannot connect")
	}

	req := &proto.ConnectRequest{
		Name:  c.Player.Name(),
		Ships: c.Player.Disposition(),
	}
	stream, err := c.bc.Connect(context.Background(), req)
	if err != nil {
		return err
	}
	log.Println("Connected, waiting for id")

	return c.Monitor(stream)
}

func (c *Client) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := c.bc.Disconnect(ctx, &proto.IdMessage{Id: c.ID})
	if err != nil {
		return err
	}

	log.Println("Disconnection succesful")
	return nil
}

func (c *Client) Monitor(stream proto.Battleship_ConnectClient) error {
	for {
		n, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error in monitoring: %v", err)
			return err
		}

		switch body := n.Body.(type) {
		// Save ID
		case *proto.Notification_ConnectReply:
			c.ID = body.ConnectReply.Id
			log.Printf("Got id %d", c.ID)
		// Handle game status
		case *proto.Notification_GameStatus:
			if body.GameStatus.Status == proto.GameStatus_VICTORY {
				log.Println("I am the king of the seas")
				return nil
			} else if body.GameStatus.Status == proto.GameStatus_DEFEAT {
				log.Println("I shall do better next time...")
				return nil
			}

			if !body.GameStatus.Play {
				continue
			}

			if err := c.Play(); err != nil {
				return err
			}
		}
	}
}

func (c *Client) Play() error {
	return c.play(0)
}

func (c *Client) play(depth int) error {
	md := 5
	if depth >= md {
		return fmt.Errorf("Maximum depth %d reached...", md)
	}

	p := c.Player.Play()
	log.Printf("Playing at position %d", p)
	req := proto.PlayRequest{
		Id:       c.ID,
		Position: p,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rep, err := c.bc.Play(ctx, &req)
	if err != nil {
		return err
	}

	if rep.Status == proto.PlayReply_NOT_YOUR_TURN {
		log.Println("Seems like it was not my turn. I'll wait")
	} else if rep.Status == proto.PlayReply_INVALID_POSITION {
		log.Println("Uh oh... I don't know how to play anymore or what?")
		return c.play(depth + 1)
	}

	c.Player.SaveResult(p, rep.Tile)
	return nil
}
