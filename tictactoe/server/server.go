package server

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/tictactoe/protos"
)

type client struct {
	Name  string
	ID    int32
	Shape Shape

	leave chan struct{}
}

type tictactoe struct {
	clients map[int32]client
	lastID  int32

	current Shape

	notifier Notifier
	game     *Game
}

func NewServer() *grpc.Server {
	s := grpc.NewServer()

	ts := tictactoe{
		clients: make(map[int32]client),
		lastID:  0,

		current: E,

		notifier: NewNotifier(),
		game:     NewGame(),
	}
	protos.RegisterTicTacToeServer(s, &ts)

	return s
}

// ----------------------------------------------------------------------------
// Connection

func (t *tictactoe) Connect(in *protos.ConnectRequest, stream protos.TicTacToe_ConnectServer) error {
	if len(t.clients) >= 2 {
		return fmt.Errorf("Server full: %d clients already connected", len(t.clients))
	}

	t.lastID += 1
	id := t.lastID

	s := X
	for _, o := range t.clients {
		if o.Shape == X {
			s = O
		}
	}

	c := client{
		Name:  in.Name,
		ID:    id,
		Shape: s,
		leave: make(chan struct{}),
	}
	t.clients[id] = c
	log.Printf("%s (%d) connected", c.Name, c.ID)

	t.notifier.Register(c.ID, stream)
	defer t.notifier.Unregister(c.ID)
	t.notifier.Notify(c.ID, t.craftConnectNotification(c))

	if len(t.clients) == 2 {
		t.startGame()
		log.Println("Let the fight begin!")
	}

	<-c.leave
	log.Printf("closing stream for %s (%d)", c.Name, c.ID)
	return nil
}

func (t *tictactoe) Disconnect(ctx context.Context, in *protos.IdMessage) (*protos.EmptyMessage, error) {
	if c, ok := t.clients[in.Id]; ok {
		close(c.leave)
	}

	log.Printf("%d disconnected", in.Id)
	delete(t.clients, in.Id)
	return &protos.EmptyMessage{}, nil
}

// ----------------------------------------------------------------------------
// Play

func (t *tictactoe) Play(ctx context.Context, in *protos.PlayRequest) (*protos.PlayReply, error) {
	c, ok := t.clients[in.Id]
	if !ok {
		return nil, fmt.Errorf("Unknown id: %d", in.Id)
	}

	if c.Shape != t.current {
		return &protos.PlayReply{
			Status: protos.PlayReply_WRONG_TURN,
		}, nil
	}

	pos := int(in.Position)
	log.Printf("%s wants to play pos %d", c.Shape, pos)

	st := t.game.Play(c.Shape, pos)

	w := t.game.Winner()
	if w != E {
		t.endGame()
		log.Printf("And the winner is: %s", w)
	} else if t.game.IsFull() {
		t.endGame()
		log.Println("It's a draw")
	} else if st == protos.PlayReply_VALID {
		t.nextPlayer()
		log.Printf("It's now to %s to play", t.current)
	}

	log.Println(t.game.Board())

	n := t.craftGameStatusNotification()
	for _, c := range t.clients {
		if err := t.notifier.Notify(c.ID, n); err != nil {
			log.Printf("Could not notify player %d: %v", c.ID, err)
		}
	}

	return &protos.PlayReply{
		Status: st,
	}, nil
}

func (t *tictactoe) startGame() {
	t.current = X

	t.game = NewGame()

	n := t.craftGameStatusNotification()
	for _, c := range t.clients {
		if err := t.notifier.Notify(c.ID, n); err != nil {
			log.Printf("Could not notify player %d: %v", c.ID, err)
		}
	}
}

func (t *tictactoe) endGame() {
	t.current = E
}

func (t *tictactoe) nextPlayer() {
	switch t.current {
	case X:
		t.current = O
	case O:
		t.current = X
	default:
		t.current = E
	}
}

// ----------------------------------------------------------------------------
// Get game status

func (t *tictactoe) GetGameStatus(ctx context.Context, in *protos.EmptyMessage) (*protos.GameStatus, error) {
	return t.gameStatus(), nil
}

// ----------------------------------------------------------------------------
// Notifications

func (t *tictactoe) craftConnectNotification(c client) *protos.Notification {
	n := &protos.ConnectReply{
		Id:    c.ID,
		Shape: int32(c.Shape),
	}
	return &protos.Notification{
		Body: &protos.Notification_ConnectReply{
			ConnectReply: n,
		},
	}
}

func (t *tictactoe) gameStatus() *protos.GameStatus {
	w := t.game.Winner()

	b := make([]int32, 9)
	for i, s := range t.game.Board() {
		b[i] = int32(s)
	}

	return &protos.GameStatus{
		Board:        b,
		CurrentShape: int32(t.current),
		Winner:       int32(w),
	}
}

func (t *tictactoe) craftGameStatusNotification() *protos.Notification {
	return &protos.Notification{
		Body: &protos.Notification_GameStatus{
			GameStatus: t.gameStatus(),
		},
	}
}
