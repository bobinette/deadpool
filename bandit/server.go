package bandit

import (
	"errors"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/bandit/components"
	"github.com/bobinette/deadpool/bandit/proto"
)

type Server struct {
	game *GameManager
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

	var n int32 = 100
	var err error
	s.game, err = NewGameManager(n)
	if err != nil {
		log.Println(err)
		return err
	}
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

	o := make(chan components.PlayerState)
	evt := components.InputEvent{
		Input:  components.Input{Arm: in.Arm},
		Output: o,
	}
	s.game.InputChannel <- evt

	var ps components.PlayerState
	select {
	case <-ctx.Done():
		return nil, errors.New("Timeout")
	case ps = <-o:
	}

	return &proto.PlayReply{
		Valid:          true,
		Score:          ps.Score,
		RemainingPlies: ps.RemainingPlies,
		Knowledge:      ps.Knowledge,
	}, nil
}
