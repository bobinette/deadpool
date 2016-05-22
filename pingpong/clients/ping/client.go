package main

import (
	"google.golang.org/grpc"

	"github.com/bobinette/deadpool/pingpong/protos"
)

type Client interface {
	Connect() error
}

type client struct {
	id  int32
	ppc protos.PingPongClient
}

func NewClient(cc *grpc.ClientConn) Client {
	return &client{
		ppc: protos.NewPingPongClient(cc),
	}
}

func (c *client) Connect() error {
	return nil
}
