package server

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	Port string
}

type Server struct {
	port string
}

func New(config Config) *Server {
	return &Server{
		port: config.Port,
	}
}

func (s *Server) Run() {
	http.HandleFunc("/connect", connectHandler)
	http.HandleFunc("/disconnect", disconnectHandler)
	http.HandleFunc("/heartbeat", heartbeatHandler)

	fmt.Println("Starting server on :" + s.port)
	if err := http.ListenAndServe(":"+s.port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Fprintf(w, "Handler 1: POST request received")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func disconnectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Fprintf(w, "Handler 2: POST request received")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Fprintf(w, "Handler 3: POST request received")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
