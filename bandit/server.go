package bandit

import (
	"errors"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/bandit/proto"
)

type Server struct {
	game *Game
	errc chan error
}

func NewServer() (*grpc.Server, <-chan error) {
	srv := Server{
		game: nil,
		errc: make(chan error),
	}

	s := grpc.NewServer()
	proto.RegisterBanditServer(s, &srv)
	return s, srv.errc
}

func (s *Server) Connect(in *proto.EmptyMessage, stream proto.Bandit_ConnectServer) error {
	if s.game != nil {
		return errors.New("Server already taken")
	}

	stream.Send(&proto.Notification{
		Msg: &proto.Notification_ConnectReply{
			ConnectReply: &proto.ConnectReply{
				NArms: 10,
			},
		},
	})

	s.game = NewGame()
	s.game.Start()
	return nil
}

func (s *Server) Disconnect(ctx context.Context, in *proto.EmptyMessage) (*proto.EmptyMessage, error) {
	s.game = nil
	return &proto.EmptyMessage{}, nil
}

func (s *Server) Play(ctx context.Context, in *proto.PlayRequest) (*proto.PlayReply, error) {
	if s.game == nil {
		return nil, errors.New("Game not started or ended prematurely")
	}

	o := make(chan float64)
	evt := InputEvent{
		Input:  Input{Arm: in.Arm},
		Output: o,
	}
	s.game.InputChannel <- evt

	var v float64
	select {
	case <-ctx.Done():
		return nil, errors.New("Timeout")
	case v = <-o:
	}

	return &proto.PlayReply{Valid: true, Value: v}, nil
}
