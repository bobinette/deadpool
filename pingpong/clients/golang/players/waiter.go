package players

import (
	"time"
)

type Waiter struct {
	d time.Duration
}

func NewWaiter() Player {
	return &Waiter{
		d: 2 * time.Second,
	}
}

func (p *Waiter) Name() string {
	return "Waiter"
}

func (p *Waiter) Play() {
	time.Sleep(p.d)
}
