package server

import (
	"context"
	"net"

	"github.com/erupshis/key_keeper/internal/server/auth"
	"github.com/erupshis/key_keeper/internal/server/sync"
	"github.com/erupshis/key_keeper/pb"
	"google.golang.org/grpc"
)

var (
	_ BaseServer = (*Server)(nil)
)

type Server struct {
	*grpc.Server
	info string
	port string
}

func NewGRPCServer(syncController *sync.Controller, authController *auth.Controller, info string, options ...grpc.ServerOption) *Server {
	s := grpc.NewServer(options...)
	pb.RegisterSyncServer(s, syncController)
	pb.RegisterAuthServer(s, authController)

	srv := &Server{
		Server: s,
		info:   info,
	}

	return srv
}

func (s *Server) Serve(lis net.Listener) error {
	return s.Server.Serve(lis)
}

func (s *Server) GracefulStop(_ context.Context) error {
	s.Server.GracefulStop()
	return nil
}
func (s *Server) GetInfo() string {
	return s.info
}

func (s *Server) Host(host string) {
	s.port = host
}
func (s *Server) GetHost() string {
	return s.port
}
