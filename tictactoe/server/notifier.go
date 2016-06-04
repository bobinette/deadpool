package server

import (
	"fmt"

	"github.com/bobinette/deadpool/tictactoe/protos"
)

type Notifier interface {
	Notify(int32, *protos.Notification) error

	Register(int32, protos.TicTacToe_ConnectServer)
	Unregister(int32)
}

type notifier struct {
	streams map[int32]protos.TicTacToe_ConnectServer
}

func NewNotifier() Notifier {
	return &notifier{
		streams: make(map[int32]protos.TicTacToe_ConnectServer),
	}
}

func (n *notifier) Notify(id int32, notif *protos.Notification) error {
	stream, ok := n.streams[id]
	if !ok {
		return fmt.Errorf("no stream for player %d", id)
	}
	return stream.Send(notif)
}

func (n *notifier) Register(id int32, stream protos.TicTacToe_ConnectServer) {
	n.streams[id] = stream
}

func (n *notifier) Unregister(id int32) {
	delete(n.streams, id)
}
