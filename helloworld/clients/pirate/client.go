package main

import (
	"fmt"
	"log"
	"math/rand"
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

	// Subscribe for notifications
	stream, err := c.hwc.Subscribe(context.Background(), &protos.IdMessage{Id: int32(c.id)})
	errCh := make(chan error)
	go func() {
		errCh <- c.Watch(stream)
	}()

	// Play at some random times
	go func() {
		src := rand.NewSource(time.Now().UnixNano())
		gen := rand.New(src)

		for {
			select {
			case <-errCh:
				return
			default:
				time.Sleep(time.Duration(gen.Intn(20)) * 250 * time.Millisecond)
				pr, err := c.hwc.Play(context.Background(), &protos.PlayRequest{Id: int32(c.id), Message: c.player.Play()})
				if err != nil {
					return
				}
				log.Println(pr.Message)
			}
		}
	}()

	// Randomly ask for notifications
	go func() {
		src := rand.NewSource(time.Now().UnixNano())
		gen := rand.New(src)

		for {
			select {
			case <-errCh:
				return
			default:
				time.Sleep(time.Duration(gen.Intn(20)) * time.Second)
				log.Println("Requesting notification, Sir!")
				n, err := c.hwc.GetNotification(context.Background(), &protos.IdMessage{Id: int32(c.id)})
				if err != nil {
					return
				}
				c.Handle(n)
			}
		}
	}()

	return <-errCh
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

func (c *client) Watch(stream protos.HelloWorld_SubscribeClient) error {
	for {
		n, err := stream.Recv()
		if err != nil {
			return err
		}
		c.Handle(n)
	}
}

func (c *client) Handle(n *protos.Notification) {
	log.Printf("In the channel: %d/%d persons", n.NPlayers, n.MaxPlayers)
	if n.Joined != "" {
		log.Printf("and %s just joined", n.Joined)
	}
	if n.Left != "" {
		log.Printf("and %s just left", n.Left)
	}
	if n.Message != "" && int(n.AuthorId) != c.id {
		log.Printf("and got message: \"%s\" from %s (id: %d)", n.Message, n.Author, n.AuthorId)
	}
}
