package grpcx

import (
	"net"

	"google.golang.org/grpc"
)

// Server represents a GRPC server.
type Server struct {
	listener net.Listener
	srv      *grpc.Server
}

// NewServer creates a new Server instance.
func NewServer(listener net.Listener, srv *grpc.Server) *Server {
	return &Server{
		listener: listener,
		srv:      srv,
	}
}

// Serve listens for incoming connections and serves RPC responses.
func (s *Server) Serve() error {
	return s.srv.Serve(s.listener)
}

// Close closes the listener and frees occupied port.
func (s *Server) Close() error {
	s.srv.GracefulStop() // It also closes the listener.

	return nil
}