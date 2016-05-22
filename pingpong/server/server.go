package server

import (
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/pingpong/protos"
)

type server struct {
	clients map[int32]string
	lastId  int32

	notifier      Notifier
	game          Game
	currentPlayer int32
}

func NewServer() *grpc.Server {
	s := grpc.NewServer()

	pps := server{
		clients: make(map[int32]string),
		lastId:  0,

		notifier:      NewNotifier(),
		game:          NewGame(),
		currentPlayer: -1,
	}
	protos.RegisterPingPongServer(s, &pps)
	return s
}

func (s *server) Connect(in *protos.ConnectRequest, stream protos.PingPong_ConnectServer) error {
	s.lastId += 1

	// Remember client
	id := s.lastId
	s.clients[id] = in.Name
	log.Printf("New client: %s (%d)", in.Name, id)

	// Register client
	s.notifier.Register(id, stream)

	s.currentPlayer = s.nextPlayer()

	// Send back it's id to the player
	n := s.idNotification(id)
	if err := s.notifier.Notify(id, n); err != nil {
		log.Printf("Error notifying %s (%d): %v", in.Name, id, err)
		return err
	}

	return nil
}

func (s *server) GetGameStatus(ctx context.Context, in *protos.IdMessage) (*protos.GameStatus, error) {
	status := s.game.Status()
	log.Printf("%d requested the game status", in.Id)
	return status, nil
}

func (s *server) Play(ctx context.Context, in *protos.PlayRequest) (*protos.PlayReply, error) {
	id := in.Id

	if id != s.currentPlayer {
		return &protos.PlayReply{Accepted: false}, nil
	}

	rep, err := s.game.Play(id)
	if err != nil {
		return nil, err
	}
	log.Printf("%d played", in.Id)
	s.currentPlayer = s.nextPlayer()
	return rep, nil
}

// ---- Notification crafting
func (s *server) idNotification(id int32) *protos.Notification {
	idMessage := &protos.IdMessage{
		Id: id,
	}
	return &protos.Notification{
		Body: &protos.Notification_IdMessage{
			IdMessage: idMessage,
		},
	}
}

// ---- Turns
func (s *server) nextPlayer() int32 {
	var ids []int32
	for k := range s.clients {
		ids = append(ids, k)
	}

	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)
	return ids[gen.Intn(len(ids))]
}
