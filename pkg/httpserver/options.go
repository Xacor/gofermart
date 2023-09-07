package httpserver

import (
	"time"
)

type Option func(*Server)

func Address(addr string) Option {
	return func(s *Server) {
		s.server.Addr = addr
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
