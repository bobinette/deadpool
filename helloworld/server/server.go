package helloworld

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/helloworld/protos"
)

type server struct {
	capacity int
	lastId   int

	players       map[int]*Player
	currentPlayer int
	playersMu     sync.Locker

	notifier Notifier
	game     Game
}

func NewServer() *grpc.Server {
	s := grpc.NewServer()
	hs := &server{
		capacity: 2,
		lastId:   0,

		players:       make(map[int]*Player),
		playersMu:     &sync.Mutex{},
		currentPlayer: -1,

		notifier: NewNotifier(),
		game:     NewGame(),
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

	s.playersMu.Lock()
	s.players[p.Id] = p
	s.lastId += 1
	s.playersMu.Unlock()

	// Notify other player there's a new sheriff in town
	for _, o := range s.players {
		go func(o *Player) {
			n := s.JoinNotification(p.Name)
			if err := s.notifier.Notify(o, n); err != nil {
				log.Println(err)
			}
		}(o)
	}

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

	s.playersMu.Lock()
	delete(s.players, p.Id)
	s.playersMu.Unlock()
	s.notifier.Unregister(p)

	// Notify other player there's man down, man down!
	for _, o := range s.players {
		go func(o *Player) {
			n := s.LeaveNotification(p.Name)
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

	if len(s.players) == s.game.Players() {
		s.startGame()

		for _, o := range s.players {
			go func(o *Player) {
				n := s.GameStatusNotification(s.currentPlayer)
				if err := s.notifier.Notify(o, n); err != nil {
					log.Println(err)
				}
			}(o)
		}
	}

	return <-ch
}

func (s *server) GetGameStatus(ctx context.Context, in *protos.IdMessage) (*protos.GameStatus, error) {
	p, ok := s.players[int(in.Id)]
	if !ok {
		return nil, fmt.Errorf("I don't know anyone corresponding to id %d", in.Id)
	}

	log.Printf("Player requested notification, id: %d", p.Id)
	return s.game.Status(), nil
}

func (s *server) Play(ctx context.Context, in *protos.PlayRequest) (*protos.PlayReply, error) {
	if int(in.Id) != s.currentPlayer {
		return &protos.PlayReply{
			Accepted: false,
			Message:  "not your turn to play.",
		}, nil
	}
	p, ok := s.players[int(in.Id)]
	if !ok {
		return nil, fmt.Errorf("I don't know anyone corresponding to id %d", in.Id)
	}

	e, err := s.game.Play(in.Message, p)
	if err != nil {
		return nil, err
	}

	s.currentPlayer = s.nextPlayer()
	for _, o := range s.players {
		go func(o *Player) {
			n := &protos.Notification{
				Event: &protos.Notification_GameEvent{
					GameEvent: e,
				},
			}
			if err := s.notifier.Notify(o, n); err != nil {
				log.Println(err)
			}
		}(o)
	}

	for _, o := range s.players {
		go func(o *Player) {
			n := s.GameStatusNotification(s.currentPlayer)
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

// ----------------------------------------------------------------------------
// Game

func (s *server) startGame() (int, error) {
	if len(s.players) != s.game.Players() {
		return -1, fmt.Errorf("incorrect number of player: %d instead of %d required", len(s.players), s.game.Players())
	}
	s.game.Reset()
	log.Println("Let the game begin!")
	return s.nextPlayer(), nil
}

func (s *server) nextPlayer() int {
	s.playersMu.Lock()
	defer s.playersMu.Unlock()

	var ids []int
	for id := range s.players {
		ids = append(ids, id)
	}

	src := rand.NewSource(time.Now().UnixNano())
	gen := rand.New(src)
	n := gen.Intn(len(ids))
	log.Printf("Next player is %d", ids[n])
	return ids[n]
}

// ----------------------------------------------------------------------------
// Server status

// Status returns the current server status
func (s *server) Status() protos.ServerStatus {
	return protos.ServerStatus{
		Max: int32(s.capacity),
		N:   int32(len(s.players)),
	}
}

// ----------------------------------------------------------------------------
// Notifications

func (s *server) JoinNotification(name string) *protos.Notification {
	st := s.Status()
	return &protos.Notification{
		Event: &protos.Notification_Join{
			Join: &protos.JoinEvent{
				PlayerName: name,
				Status:     &st,
			},
		},
	}
}

func (s *server) LeaveNotification(name string) *protos.Notification {
	st := s.Status()
	return &protos.Notification{
		Event: &protos.Notification_Leave{
			&protos.LeaveEvent{
				PlayerName: name,
				Status:     &st,
			},
		},
	}
}

func (s *server) GameStatusNotification(currentPlayer int) *protos.Notification {
	st := s.game.Status()
	st.CurrentPlayer = int32(currentPlayer)
	return &protos.Notification{
		Event: &protos.Notification_GameStatus{
			st,
		},
	}
}
