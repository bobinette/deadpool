package pingpong

import (
	"log"
)

type Event int

const (
	EventConnect    Event = iota
	EventDisconnect Event = iota
	EventGameStart  Event = iota
	EventGameOver   Event = iota
	EventPly        Event = iota
)

type Logger interface {
	Log(e Event, args ...interface{})

	Error(string)
	Errorf(format string, args ...interface{})
}

type TerminalLogger struct{}

func (l *TerminalLogger) Log(e Event, args ...interface{}) {
	switch e {
	case EventConnect:
		log.Printf("%s just joined the game", args[0].(*Player).Name)
	case EventDisconnect:
		log.Printf("%s left", args[0].(*Player).Name)
	case EventGameStart:
		log.Printf("Game starts: %d plies", args[0].(int))
	case EventGameOver:
		log.Printf("Game over after %d plies", args[0].(int))
	case EventPly:
		p := args[0].(*Player)
		log.Printf("%d) %s: %s", args[1].(int), p.Name, p.Sound.String())
	}
}
func (l *TerminalLogger) Error(s string) {
	log.Println(s)
}

func (l *TerminalLogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
