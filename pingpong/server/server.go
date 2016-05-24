package server

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/pingpong/protos"
)

type server struct {
	// Clients
	clients map[int32]client
	// Id given last to a client. To ensure id uniqueness
	lastId int32
	// Id of the client that can play
	currentSound Sound

	notifier Notifier
	game     Game
}

type client struct {
	*Player

	leave chan struct{}
}

func NewServer() *grpc.Server {
	s := grpc.NewServer()

	pps := server{
		clients:      make(map[int32]client),
		lastId:       0,
		currentSound: Silence,

		notifier: NewNotifier(),
		game:     NewGame(),
	}
	protos.RegisterPingPongServer(s, &pps)
	return s
}

func (s *server) Connect(in *protos.ConnectRequest, stream protos.PingPong_ConnectServer) error {
	// Remember client
	c, err := s.newClient()
	if err != nil {
		return err
	}
	s.clients[c.Id] = c
	log.Printf("New client: %s (%d)", in.Name, c.Id)

	// Register client
	s.notifier.Register(c.Id, stream)

	// Send back it's id to the player
	n := s.craftConnectReply(c.Id, c.Sound)
	if err := s.notifier.Notify(c.Id, n); err != nil {
		log.Printf("Error notifying %s (%d): %v", in.Name, c.Id, err)
		return err
	}

	if len(s.clients) >= 2 {
		log.Println("Starting game...")
		s.currentSound = s.nextSound()
		for _, cl := range s.clients {
			go func(cl client) {
				s.notifier.Notify(cl.Id, s.craftGameStatusNotification())
			}(cl)
		}
	}

	<-c.leave
	return nil
}

func (s *server) Leave(ctx context.Context, in *protos.IdMessage) (*protos.LeaveReply, error) {
	id := in.Id
	s.notifier.Unregister(id)
	if c, ok := s.clients[id]; ok {
		close(c.leave)
	}
	delete(s.clients, id)
	log.Printf("Client %d left", id)
	return &protos.LeaveReply{}, nil
}

func (s *server) GetGameStatus(ctx context.Context, in *protos.IdMessage) (*protos.GameStatus, error) {
	status := s.craftGameStatus()
	// log.Printf("%d requested the game status", in.Id)
	return status, nil
}

func (s *server) Play(ctx context.Context, in *protos.PlayRequest) (*protos.PlayReply, error) {
	id := in.Id

	c, ok := s.clients[id]
	if !ok {
		return nil, fmt.Errorf("Id %d unknown", id)
	}

	if c.Sound != s.currentSound {
		return &protos.PlayReply{Accepted: false}, nil
	}

	rep, err := s.game.Play(c.Sound)
	if err != nil {
		return nil, err
	}
	log.Println(c.Sound)
	s.currentSound = s.nextSound()

	for _, cl := range s.clients {
		go func(cl client) {
			s.notifier.Notify(cl.Id, s.craftGameStatusNotification())
		}(cl)
	}

	return rep, nil
}

// ---- Notification crafting
func (s *server) craftConnectReply(id int32, sound Sound) *protos.Notification {
	n := &protos.ConnectReply{
		Id:    id,
		Sound: int32(sound),
	}
	return &protos.Notification{
		Body: &protos.Notification_ConnectReply{
			ConnectReply: n,
		},
	}
}

func (s *server) craftGameStatus() *protos.GameStatus {
	return &protos.GameStatus{
		Pingpong:     s.game.History(),
		CurrentSound: int32(s.currentSound),
	}
}

func (s *server) craftGameStatusNotification() *protos.Notification {
	return &protos.Notification{
		Body: &protos.Notification_GameStatus{
			GameStatus: s.craftGameStatus(),
		},
	}
}

// ---- Client management
func (s *server) newClient() (client, error) {
	if len(s.clients) >= 2 {
		return client{}, fmt.Errorf("Server full (%d client already)", len(s.clients))
	}

	s.lastId += 1
	id := s.lastId
	sound := Ping

	for _, c := range s.clients {
		if c.Sound == Ping {
			sound = Pong
		}
	}

	c := client{
		&Player{
			Id:    id,
			Sound: sound,
		},
		make(chan struct{}),
	}
	return c, nil
}

func (s *server) nextSound() Sound {
	switch s.currentSound {
	case Ping:
		return Pong
	case Pong:
		return Ping
	}

	sounds := []Sound{Ping, Pong}
	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)
	return sounds[gen.Intn(len(sounds))]
}
