package main

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/tictactoe/protos"
)

type Client interface {
	Connect() error
	Disconnect() error
}

const name = "Random"

type Shape int32

const (
	E Shape = iota // Empty
	X Shape = iota
	O Shape = iota
)

func (s Shape) String() string {
	switch s {
	case X:
		return "X"
	case O:
		return "O"
	}
	return "-"
}

type client struct {
	id    int32
	shape Shape
	board []Shape

	gen *rand.Rand

	tc protos.TicTacToeClient
}

func NewClient(cc *grpc.ClientConn) Client {
	src := rand.NewSource(time.Now().UnixNano())
	return &client{
		board: make([]Shape, 9),
		gen:   rand.New(src),
		tc:    protos.NewTicTacToeClient(cc),
	}
}

func (c *client) Connect() error {
	stream, err := c.tc.Connect(context.Background(), &protos.ConnectRequest{Name: name})
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
	c.id = rep.Id
	c.shape = Shape(rep.Shape)
	log.Printf("Got id %d and shape %s", c.id, c.shape)
	return c.Monitor(stream)
}

func (c *client) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := c.tc.Disconnect(ctx, &protos.IdMessage{Id: c.id})
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Monitor(stream protos.TicTacToe_ConnectClient) error {
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

func (c *client) Handle(n *protos.Notification) error {
	switch body := n.Body.(type) {
	case *protos.Notification_GameStatus:
		for i, s := range body.GameStatus.Board {
			c.board[i] = Shape(s)
		}
		log.Println(c.board)

		if Shape(body.GameStatus.Winner) == c.shape {
			log.Println("I won!")
		} else if Shape(body.GameStatus.Winner) != E {
			log.Println("I lost...")
		} else if Shape(body.GameStatus.CurrentShape) == c.shape {
			return c.Play()
		}
	}
	return nil
}

func (c *client) Play() error {
	ps := make([]int, 0, 9)

	for p, s := range c.board {
		if Shape(s) == E {
			ps = append(ps, p)
		}
	}
	if len(ps) == 0 {
		return errors.New("No available position to play")
	}

	p := ps[c.gen.Intn(len(ps))]
	req := &protos.PlayRequest{
		Id:       c.id,
		Position: int32(p),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	log.Printf("Playing at position %d", p)
	rep, err := c.tc.Play(ctx, req)
	if err != nil {
		return err
	}

	switch rep.Status {
	case protos.PlayReply_VALID:
		log.Println("Accepted")
	case protos.PlayReply_WRONG_TURN:
		log.Println("Not my turn to play...")
	case protos.PlayReply_POSITION_TAKEN:
		log.Println("Position taken")
	case protos.PlayReply_INVALID_POSITION:
		log.Println("Invalid position")
	}
	return nil
}
