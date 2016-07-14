package pipeline

type Channel struct {
	Key string
	Ch  chan interface{}
}

func NewChannel(key string) *Channel {
	return &Channel{
		Key: key,
		Ch:  make(chan interface{}),
	}
}
