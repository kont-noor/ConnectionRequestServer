package server

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Config struct {
	Port   string
	Router http.Handler
	Log    *zap.Logger
}

type Server struct {
	port       string
	httpServer *http.Server
	log        *zap.Logger
}

func New(config Config) *Server {
	return &Server{
		port: config.Port,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("localhost:%s", config.Port),
			Handler: config.Router,
		},
		log: config.Log,
	}
}

func (s *Server) Run() {
	s.log.Info("Starting server on :" + s.port)
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		s.log.Sugar().Errorf("Could not start server: %v\n", err)
	}
}
