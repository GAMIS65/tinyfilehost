package server

import (
	"log/slog"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	logger *slog.Logger
}

func New(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr: addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}