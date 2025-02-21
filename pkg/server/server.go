package server

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Config struct {
	Hostname string
	Port     string
	Router   http.Handler
	Log      *zap.Logger
}

type Server struct {
	hostname   string
	port       string
	httpServer *http.Server
	log        *zap.Logger
}

func New(config Config) *Server {
	return &Server{
		port:     config.Port,
		hostname: config.Hostname,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", config.Hostname, config.Port),
			Handler: config.Router,
		},
		log: config.Log,
	}
}

func (s *Server) Run() {
	s.log.Sugar().Infof("Starting server on %s:%s", s.hostname, s.port)
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		s.log.Sugar().Errorf("Could not start server: %v\n", err)
	}
}
