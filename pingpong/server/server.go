package server

import (
	"google.golang.org/grpc"
)

type server struct {
}

func NewServer() *grpc.Server {
	s := grpc.NewServer()
	return s
}
