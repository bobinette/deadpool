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
			if !body.GameStatus.Play {
				continue
			}
			if err := c.Play(); err != nil {
				return err
			}
		}
	}
}

var i = 3

func (c *Client) Play() error {
	i -= 1
	if i < 0 {
		return fmt.Errorf("Finished")
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

	log.Printf("Status: %s", proto.Tile_name[int32(rep.Status)])
	return nil
}

type Board []proto.Tile

func (b Board) String() string {
	s := "\n"
	s += "------------\n"
	for i := 0; i < 10; i++ {
		s += "|"
		for j := 0; j < 10; j++ {
			switch b[10*i+j] {
			case proto.Tile_SEA:
				s += " "
			case proto.Tile_UNKNOWN:
				s += "-"
			case proto.Tile_SHIP:
				s += "O"
			case proto.Tile_SUNKEN:
				s += "X"
			}
		}
		s += "|"
		s += "\n"
	}
	s += "------------"
	return s
}
