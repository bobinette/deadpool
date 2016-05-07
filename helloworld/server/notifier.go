package helloworld

import (
	"fmt"
	"log"

	"github.com/bobinette/deadpool/helloworld/protos"
)

type Notifier interface {
	Notify(*Player, *protos.Notification) error

	Register(*Player, protos.HelloWorld_SubscribeServer) chan error
	Unregister(*Player)
}

type stream struct {
	protos.HelloWorld_SubscribeServer

	errChan chan error
}

type notifier struct {
	streams map[int]stream
}

func NewNotifier() Notifier {
	return &notifier{
		streams: make(map[int]stream),
	}
}

func (n *notifier) Register(p *Player, s protos.HelloWorld_SubscribeServer) chan error {
	ch := make(chan error)
	n.streams[p.Id] = stream{s, ch}
	return ch
}

func (n *notifier) Unregister(p *Player) {
	if stream, ok := n.streams[p.Id]; ok {
		close(stream.errChan)
	}
	delete(n.streams, p.Id)
}

func (n *notifier) Notify(p *Player, r *protos.Notification) error {
	stream, ok := n.streams[p.Id]
	if !ok {
		return fmt.Errorf("no stream for player %d", p.Id)
	}

	err := stream.Send(r)
	if err != nil {
		go func() {
			stream.errChan <- err
		}()
	}
	log.Printf("Sent notification to player %d", p.Id)
	return err
}
