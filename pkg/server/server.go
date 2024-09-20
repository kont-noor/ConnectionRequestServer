package server

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	Port   string
	Router http.Handler
}

type Server struct {
	port       string
	httpServer *http.Server
}

func New(config Config) *Server {
	return &Server{
		port: config.Port,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("localhost:%s", config.Port),
			Handler: config.Router,
		},
	}
}

func (s *Server) Run() {
	fmt.Println("Starting server on :" + s.port)
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
