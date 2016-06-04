package battleship

import (
	"fmt"

	"github.com/bobinette/deadpool/battleship/proto"
)

type Notifier interface {
	Notify(*Client, *proto.Notification) error

	Register(*Client, proto.Battleship_ConnectServer)
	Unregister(*Client)
}

type notifier struct {
	streams map[int32]proto.Battleship_ConnectServer
}

func NewNotifier() Notifier {
	return &notifier{
		streams: make(map[int32]proto.Battleship_ConnectServer),
	}
}

func (nf *notifier) Notify(c *Client, n *proto.Notification) error {
	s, ok := nf.streams[c.ID]
	if !ok {
		return fmt.Errorf("No stream for client %d", c.ID)
	}

	return s.Send(n)
}

func (nf *notifier) Register(c *Client, stream proto.Battleship_ConnectServer) {
	nf.streams[c.ID] = stream
}

func (nf *notifier) Unregister(c *Client) {
	delete(nf.streams, c.ID)
}
