package helloworld

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/helloworld/protos"
)

type server struct {
	capacity int
	players  map[int]*Player
	lastId   int

	notifier Notifier
}

func NewServer() *grpc.Server {
	s := grpc.NewServer()
	hs := &server{
		capacity: 2,
		players:  make(map[int]*Player),
		lastId:   0,

		notifier: NewNotifier(),
	}
	protos.RegisterHelloWorldServer(s, hs)

	return s
}

func (s *server) Connect(ctx context.Context, in *protos.ConnectRequest) (*protos.IdMessage, error) {
	if len(s.players) >= s.capacity {
		return nil, fmt.Errorf("Sorry, too many people here (capacity %d)", s.capacity)
	}

	if in.Name == "" {
		return nil, fmt.Errorf("Nah, you need a name!")
	}

	p := NewPlayer(s.lastId, in.Name)
	s.lastId += 1

	// Notify other player there's a new sheriff in town
	for _, o := range s.players {
		go func(o *Player) {
			n := s.NotificationForPlayer(p)
			n.Joined = p.Name
			if err := s.notifier.Notify(o, n); err != nil {
				log.Println(err)
			}
		}(o)
	}

	s.players[p.Id] = p
	log.Printf("New player, id: %d", p.Id)
	return &protos.IdMessage{
		Id: int32(p.Id),
	}, nil
}

func (s *server) Leave(ctx context.Context, in *protos.IdMessage) (*protos.LeaveReply, error) {
	p, ok := s.players[int(in.Id)]
	if !ok {
		return nil, fmt.Errorf("I don't know anyone corresponding to id %d", in.Id)
	}

	delete(s.players, p.Id)
	s.notifier.Unregister(p)

	// Notify other player there's man down, man down!
	for _, o := range s.players {
		go func(o *Player) {
			n := s.NotificationForPlayer(o)
			n.Left = p.Name
			if err := s.notifier.Notify(o, n); err != nil {
				log.Println(err)
			}
		}(o)
	}

	log.Printf("Player left, id: %d", p.Id)
	return &protos.LeaveReply{
		Time: p.ConnectedFor().Nanoseconds() / int64(time.Millisecond),
	}, nil
}

func (s *server) Subscribe(in *protos.IdMessage, stream protos.HelloWorld_SubscribeServer) error {
	p, ok := s.players[int(in.Id)]
	if !ok {
		return fmt.Errorf("I don't know anyone corresponding to id %d", in.Id)
	}

	// Register the new player to the notifier, and defer unregistration
	log.Printf("Player subscribed, id: %d", p.Id)
	ch := s.notifier.Register(p, stream)
	defer s.notifier.Unregister(p)

	return <-ch
}

func (s *server) GetNotification(ctx context.Context, in *protos.IdMessage) (*protos.Notification, error) {
	p, ok := s.players[int(in.Id)]
	if !ok {
		return nil, fmt.Errorf("I don't know anyone corresponding to id %d", in.Id)
	}

	log.Printf("Player requested notification, id: %d", p.Id)
	return s.NotificationForPlayer(p), nil
}

func (s *server) Play(ctx context.Context, in *protos.PlayRequest) (*protos.PlayReply, error) {
	p, ok := s.players[int(in.Id)]
	if !ok {
		return nil, fmt.Errorf("I don't know anyone corresponding to id %d", in.Id)
	}

	for _, o := range s.players {
		go func(o *Player) {
			n := s.NotificationForPlayer(o)
			n.Message = in.Message
			n.Author = p.Name
			n.AuthorId = int32(p.Id)
			if err := s.notifier.Notify(o, n); err != nil {
				log.Println(err)
			}
		}(o)
	}

	return &protos.PlayReply{
		Accepted: true,
		Message:  in.Message,
	}, nil
}

func (s *server) NotificationForPlayer(p *Player) *protos.Notification {
	return &protos.Notification{
		MaxPlayers: int32(s.capacity),
		NPlayers:   int32(len(s.players)),
		Joined:     "",
		Left:       "",
	}
}
