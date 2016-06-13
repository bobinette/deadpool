package battleship

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/battleship/proto"
)

type Client struct {
	Name   string
	ID     int32
	Placed bool

	leave chan struct{}
}

type Server struct {
	clients map[int32]*Client
	lastID  int32
	current int32
	locker  sync.Locker

	notifier Notifier

	maxGames int // Maximum number of games to play
	nGames   int // Number of games already played
	game     *Game
	plies    int

	stopCh chan error
}

var maxGames = 1

func init() {
	flag.IntVar(&maxGames, "games", maxGames, "The number of games to be played in a row")
}

func NewServer() (*grpc.Server, <-chan error) {
	s := grpc.NewServer()

	if maxGames <= 0 {
		maxGames = 1
	}

	srv := Server{
		clients: make(map[int32]*Client, 2),
		lastID:  0,
		current: 0,
		locker:  &sync.Mutex{},

		notifier: NewNotifier(),

		maxGames: maxGames,
		nGames:   0,
		game:     nil,
		plies:    0,

		stopCh: make(chan error, 1),
	}
	proto.RegisterBattleshipServer(s, &srv)
	log.Printf("Will play %d games", srv.maxGames)

	ch := make(chan error, 1)
	go func() {
		err := <-srv.stopCh
		for _, c := range srv.clients {
			c.leave <- struct{}{}
		}
		ch <- err
	}()

	return s, ch
}

func (s *Server) Connect(in *proto.ConnectRequest, stream proto.Battleship_ConnectServer) error {
	if len(s.clients) >= 2 {
		return fmt.Errorf("2 clients already connected")
	}

	s.lastID += 1
	c := Client{
		ID:     s.lastID,
		Name:   in.Name,
		Placed: false,
		leave:  make(chan struct{}),
	}

	s.clients[c.ID] = &c
	defer delete(s.clients, c.ID)

	s.notifier.Register(&c, stream)
	defer s.notifier.Unregister(&c)
	if err := s.notifier.Notify(&c, s.ConnectReplyNotification(&c)); err != nil {
		return err
	}

	if len(s.clients) == 2 {
		s.game = NewGame()
		s.DispatchGameWillStartNotifications()
	}

	<-c.leave
	return nil
}

func (s *Server) Disconnect(ctx context.Context, in *proto.IdMessage) (*proto.EmptyMessage, error) {
	c, ok := s.clients[in.Id]
	if !ok {
		return nil, fmt.Errorf("Unknown id %d", in.Id)
	}

	s.locker.Lock()
	s.current = -1
	s.locker.Unlock()
	close(c.leave)
	return &proto.EmptyMessage{}, nil

}

func (s *Server) Play(ctx context.Context, in *proto.PlayRequest) (*proto.PlayReply, error) {
	c, ok := s.clients[in.Id]
	if !ok {
		return nil, fmt.Errorf("Unknown id %d", in.Id)
	}

	if in.Id != s.current {
		return &proto.PlayReply{Tile: proto.Tile_UNKNOWN, Status: proto.PlayReply_NOT_YOUR_TURN}, nil
	}

	if in.Position < 0 || in.Position > 100 {
		return &proto.PlayReply{Tile: proto.Tile_UNKNOWN, Status: proto.PlayReply_INVALID_POSITION}, nil
	}

	t := s.game.RegisterPly(c.ID, in.Position)
	if t == proto.Tile_SHIP {
		log.Printf("Player %d touched at %d", c.ID, in.Position)
	} else if t == proto.Tile_SUNK {
		log.Printf("Player %d sank a ship at %d", c.ID, in.Position)
	}

	s.locker.Lock()
	defer s.locker.Unlock()
	for id := range s.clients {
		if id != c.ID {
			s.current = id
		}
	}

	if err := s.DispatchGameStatusNotifications(); err != nil {
		log.Printf("error dispatching game status notification: %v", err)
	}

	s.plies += 1
	w := s.game.Winner()
	for id, c := range s.clients {
		if w == id {
			log.Printf("%s (%d) won after %d plies", c.Name, c.ID, s.plies)
		}
	}
	if w != -1 {
		s.game = nil
		s.nGames += 1
		if s.nGames < s.maxGames {
			s.game = NewGame()
			if err := s.DispatchGameWillStartNotifications(); err != nil {
				return nil, err
			}
		} else {
			for _, c := range s.clients {
				c.leave <- struct{}{}
			}
			s.stopCh <- nil
		}
	}

	return &proto.PlayReply{Tile: t, Status: proto.PlayReply_ACCEPTED}, nil
}

func (s *Server) Place(ctx context.Context, in *proto.PlaceRequest) (*proto.PlaceReply, error) {
	c, ok := s.clients[in.Id]
	if !ok {
		return nil, fmt.Errorf("Unknown id %d", c.ID)
	}

	err := s.game.SavePlacement(c.ID, in.Ships)
	if err != nil {
		log.Printf("Error in ship placement: %v", err)
		return &proto.PlaceReply{Valid: false}, nil
	}

	c.Placed = true
	canStart := true
	for _, c := range s.clients {
		canStart = canStart && c.Placed
	}
	if canStart {
		s.locker.Lock()
		s.plies = 0
		s.current = c.ID
		s.locker.Unlock()
		log.Printf("Let's begin with %d", s.current)
		if err := s.DispatchGameStatusNotifications(); err != nil {
			log.Printf("error dispatching game status notification: %v", err)
		}
	}

	return &proto.PlaceReply{Valid: true}, nil
}

// ----------------------------------------------------------------------------
// Send notifications

func (s *Server) DispatchGameStatusNotifications() error {
	for _, c := range s.clients {
		if err := s.notifier.Notify(c, s.GameStatusNotification(c)); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) DispatchGameWillStartNotifications() error {
	for _, c := range s.clients {
		if err := s.notifier.Notify(c, s.GameWillStartNotification()); err != nil {
			return err
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// Craft notifications

func (s *Server) ConnectReplyNotification(c *Client) *proto.Notification {
	return &proto.Notification{
		Body: &proto.Notification_ConnectReply{
			ConnectReply: &proto.IdMessage{Id: c.ID},
		},
	}
}

func (s *Server) GameStatusNotification(c *Client) *proto.Notification {
	play := s.current == c.ID
	status := proto.GameStatus_PLAYING
	w := s.game.Winner()
	if c.ID == w {
		status = proto.GameStatus_VICTORY
	} else if w != -1 {
		status = proto.GameStatus_DEFEAT

	}
	return &proto.Notification{
		Body: &proto.Notification_GameStatus{
			GameStatus: &proto.GameStatus{
				Play:   play,
				Status: status,
			},
		},
	}
}

func (s *Server) GameWillStartNotification() *proto.Notification {
	return &proto.Notification{
		Body: &proto.Notification_GameWillStart{
			GameWillStart: &proto.EmptyMessage{},
		},
	}
}
