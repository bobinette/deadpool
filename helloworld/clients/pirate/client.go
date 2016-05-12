package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/golang/protobuf/ptypes"
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

	// Randomly ask for notifications
	// go func() {
	// 	src := rand.NewSource(time.Now().UnixNano())
	// 	gen := rand.New(src)

	// 	for {
	// 		select {
	// 		case <-errCh:
	// 			return
	// 		default:
	// 			time.Sleep(time.Duration(gen.Intn(20)) * time.Second)
	// 			log.Println("Requesting notification, Sir!")
	// 			n, err := c.hwc.GetGameStatus(context.Background(), &protos.IdMessage{Id: int32(c.id)})
	// 			if err != nil {
	// 				return
	// 			}
	// 			log.Printf("%d messages in the chat", len(n.Chat))
	// 		}
	// 	}
	// }()

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

		if err := c.Handle(n); err != nil {
			return err
		}
	}
}

func (c *client) Handle(n *protos.Notification) error {
	switch e := n.GetEvent().(type) {
	case *protos.Notification_Join:
		name := e.Join.PlayerName
		status := e.Join.Status
		log.Printf("%s just joined! In the channel: %d/%d persons", name, status.N, status.Max)
	case *protos.Notification_GameEvent:
		t, err := ptypes.Timestamp(e.GameEvent.Timestamp)
		if err != nil {
			return err
		}
		log.Printf("%s - %s says:\n\t%s", t.String(), e.GameEvent.Author, e.GameEvent.Message)
	case *protos.Notification_GameStatus:
		if int(e.GameStatus.CurrentPlayer) != c.id {
			return nil
		}
		log.Println("My turn to play!!")
		pr, err := c.hwc.Play(context.Background(), c.Play())
		if err != nil {
			return err
		}
		if !pr.Accepted {

			log.Printf("Play not accepted: %s", pr.Message)
		} else {
			log.Println("Message accepted. Sent: %s", pr.Message)
		}
	case nil:
		log.Println("Nil notification received")
	default:
		log.Println("Unknown notification type")
	}
	return nil
}

func (c *client) Play() *protos.PlayRequest {
	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)
	time.Sleep(time.Duration(gen.Intn(20)) * time.Second)
	return &protos.PlayRequest{Id: int32(c.id), Message: c.player.Play()}
}
