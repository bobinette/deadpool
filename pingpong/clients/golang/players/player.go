package players

type Player interface {
	// Name returns the name of the player. This is used for tracking which
	// AI is playing, so clear names are recommented
	Name() string

	// Play is used to pause the client before it should fire the sound
	Play()
}

func NewPlayer(name string) Player {
	switch name {
	case "waiter":
		return NewWaiter()
	}
	return nil
}
