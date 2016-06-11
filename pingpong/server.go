package pingpong

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/pingpong/proto"
)

type server struct {
	// Clients
	clients map[int32]client
	// Id given last to a client. To ensure id uniqueness
	lastId int32
	// Id of the client that can play
	currentSound Sound

	maxPlies int
	plies    int

	notifier Notifier
	game     Game
	logger   Logger
}

type client struct {
	*Player

	leave chan struct{}
}

func NewServer() *grpc.Server {
	s := grpc.NewServer()

	logger := TerminalLogger{}
	pps := server{
		clients:      make(map[int32]client),
		lastId:       0,
		currentSound: Silence,

		maxPlies: 0,
		plies:    0,

		notifier: NewNotifier(),
		game:     NewGame(),
		logger:   &logger,
	}
	proto.RegisterPingPongServer(s, &pps)
	return s
}

func (s *server) Connect(in *proto.ConnectRequest, stream proto.PingPong_ConnectServer) error {
	// Remember client
	c, err := s.newClient(in.Name)
	if err != nil {
		return err
	}
	s.clients[c.Id] = c
	s.logger.Log(EventConnect, c.Player)

	// Register client
	s.notifier.Register(c.Id, stream)

	// Send back it's id to the player
	n := s.craftConnectReply(c.Id, c.Sound)
	if err := s.notifier.Notify(c.Id, n); err != nil {
		s.logger.Errorf("Error notifying %s (%d): %v", in.Name, c.Id, err)
		return err
	}

	if len(s.clients) >= 2 {
		src := rand.NewSource(time.Now().UnixNano())
		gen := rand.New(src)

		s.maxPlies = gen.Intn(81) + 20
		s.logger.Log(EventGameStart, s.maxPlies)

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

func (s *server) Disconnect(ctx context.Context, in *proto.IdMessage) (*proto.DisconnectReply, error) {
	id := in.Id
	s.notifier.Unregister(id)
	if c, ok := s.clients[id]; ok {
		close(c.leave)
		s.logger.Log(EventDisconnect, c.Player)
	}
	delete(s.clients, id)
	return &proto.DisconnectReply{}, nil
}

func (s *server) GetGameStatus(ctx context.Context, in *proto.IdMessage) (*proto.GameStatus, error) {
	status := s.craftGameStatus()
	// log.Printf("%d requested the game status", in.Id)
	return status, nil
}

func (s *server) Play(ctx context.Context, in *proto.PlayRequest) (*proto.PlayReply, error) {
	id := in.Id

	c, ok := s.clients[id]
	if !ok {
		return nil, fmt.Errorf("Id %d unknown", id)
	}

	if c.Sound != s.currentSound {
		return &proto.PlayReply{Accepted: false}, nil
	}

	rep, err := s.game.Play(c.Sound)
	if err != nil {
		return nil, err
	}
	s.plies += 1
	s.logger.Log(EventPly, c.Player, s.plies)

	if s.plies >= s.maxPlies {
		for _, c := range s.clients {
			close(c.leave)
		}
		s.logger.Log(EventGameOver, s.plies)
	} else {
		s.currentSound = s.nextSound()

		for _, cl := range s.clients {
			go func(cl client) {
				s.notifier.Notify(cl.Id, s.craftGameStatusNotification())
			}(cl)
		}
	}

	return rep, nil
}

// ---- Notification crafting
func (s *server) craftConnectReply(id int32, sound Sound) *proto.Notification {
	n := &proto.ConnectReply{
		Id:    id,
		Sound: int32(sound),
	}
	return &proto.Notification{
		Body: &proto.Notification_ConnectReply{
			ConnectReply: n,
		},
	}
}

func (s *server) craftGameStatus() *proto.GameStatus {
	return &proto.GameStatus{
		Pingpong:     s.game.History(),
		CurrentSound: int32(s.currentSound),
	}
}

func (s *server) craftGameStatusNotification() *proto.Notification {
	return &proto.Notification{
		Body: &proto.Notification_GameStatus{
			GameStatus: s.craftGameStatus(),
		},
	}
}

// ---- Client management
func (s *server) newClient(name string) (client, error) {
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
			Name:  name,
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
