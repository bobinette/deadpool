package players

import (
	"github.com/bobinette/deadpool/bandit/components"
)

type Player interface {
	Play() int32
	Save(components.PlayerState)
}
