package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/helloworld/protos"
)

type Client interface {
	Connect() error
	Leave() error
}

type client struct {
	id     int
	player Player
	hwc    protos.HelloWorldClient
}

func NewClient(cc *grpc.ClientConn) Client {
	hwc := protos.NewHelloWorldClient(cc)

	return &client{
		player: NewPlayer(),
		hwc:    hwc,
	}
}

func (c *client) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rep, err := c.hwc.Connect(ctx, &protos.ConnectRequest{Name: c.player.Name()})
	if err != nil {
		return fmt.Errorf("could not connect: %v", err)
	}

	c.id = int(rep.Id)
	log.Printf("Connected with id %d", c.id)

	pr, err := c.hwc.Play(context.Background(), &protos.PlayRequest{Id: int32(c.id), Message: c.player.Play()})
	if err != nil {
		return err
	}
	log.Println(pr.Message)

	return nil
}

func (c *client) Leave() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rep, err := c.hwc.Leave(ctx, &protos.IdMessage{Id: int32(c.id)})
	if err != nil {
		return fmt.Errorf("could not connect: %v", err)
	}

	log.Printf("Time connected: %dms", rep.Time)
	return nil
}
