package server

import (
	"connection_request_server/pkg/mongo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	Port  string
	Mongo *mongo.Mongo
}

type Server struct {
	port  string
	mongo *mongo.Mongo
}

type requestParams struct {
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
}

func New(config Config) *Server {
	return &Server{
		port:  config.Port,
		mongo: config.Mongo,
	}
}

func (s *Server) Run() {
	http.HandleFunc("/connect", s.connectHandler)
	http.HandleFunc("/disconnect", s.disconnectHandler)
	http.HandleFunc("/heartbeat", s.heartbeatHandler)

	fmt.Println("Starting server on :" + s.port)
	if err := http.ListenAndServe(":"+s.port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

func (s *Server) connectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	fmt.Println(r.Method)
	if r.Method == http.MethodPost {
		var params requestParams
		if err := parseRequest(r, &params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		connection, err := s.mongo.FindActiveConnection(params.UserID, params.DeviceID)
		// TODO: this also rises the error if connection is not found; need to fix this
		//if err != nil {
		//	http.Error(w, "Failed to find connection"+err.Error(), http.StatusInternalServerError)
		//	return
		//}

		if connection != nil {
			if connection.LastHeartbeat.Time().Before(time.Now().Add(-5 * time.Second)) {
				fmt.Fprintf(w, "Connection expired; User ID: %s, Device ID: %s, Last heartbeat: %s", params.UserID, params.DeviceID, connection.LastHeartbeat.Time().String())
				err = s.mongo.DeleteConnection(params.UserID, params.DeviceID)
				if err != nil {
					http.Error(w, "Failed to delete connection", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Connection already exists", http.StatusConflict)
				return
			}
		}

		newConnection := mongo.Connection{
			ID:            primitive.NewObjectID(),
			UserID:        primitive.Symbol(params.UserID),
			DeviceID:      primitive.Symbol(params.DeviceID),
			ConnectedAt:   primitive.NewDateTimeFromTime(time.Now()),
			LastHeartbeat: primitive.NewDateTimeFromTime(time.Now()),
		}

		err = s.mongo.InsertConnection(newConnection)
		if err != nil {
			http.Error(w, "Failed to insert connection"+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Connection approved; User ID: %s, Device ID: %s", params.UserID, params.DeviceID)
	} else {
		fmt.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func parseRequest(r *http.Request, params *requestParams) error {
	return json.NewDecoder(r.Body).Decode(params)
}

func (s *Server) disconnectHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	fmt.Println(r.Method)
	if r.Method == http.MethodPost {
		var params requestParams
		if err := parseRequest(r, &params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		connection, err := s.mongo.FindActiveConnection(params.UserID, params.DeviceID)
		if err != nil {
			http.Error(w, "Failed to find connection"+err.Error(), http.StatusInternalServerError)
			return
		} else if connection == nil {
			http.Error(w, "Connection not found", http.StatusNotFound)
			return
		} else {
			err = s.mongo.DeleteConnection(params.UserID, params.DeviceID)
			if err != nil {
				http.Error(w, "Failed to delete connection"+err.Error(), http.StatusInternalServerError)
				return
			}

			fmt.Fprintf(w, "Connection deleted; User ID: %s, Device ID: %s", params.UserID, params.DeviceID)
		}
	} else {
		fmt.Println("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (s *Server) heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var params requestParams
		if err := parseRequest(r, &params); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		connection, err := s.mongo.FindActiveConnection(params.UserID, params.DeviceID)
		if err != nil {
			http.Error(w, "Failed to find connection"+err.Error(), http.StatusInternalServerError)
			return
		} else if connection == nil {
			http.Error(w, "Connection not found", http.StatusNotFound)
			return
		} else {
			if connection.LastHeartbeat.Time().Before(time.Now().Add(-5 * time.Second)) {
				http.Error(w, "Connection expired", http.StatusConflict)
				err = s.mongo.DeleteConnection(params.UserID, params.DeviceID)
				if err != nil {
					http.Error(w, "Failed to delete connection", http.StatusInternalServerError)
					return
				}
			} else {
				err = s.mongo.HeartbeatConnection(params.UserID, params.DeviceID)
				if err != nil {
					http.Error(w, "Failed to update connection"+err.Error(), http.StatusInternalServerError)
					return
				}

				fmt.Fprintf(w, "Heartbeat received; User ID: %s, Device ID: %s", params.UserID, params.DeviceID)
			}
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
