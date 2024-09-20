package service

import (
	"connection_request_server/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Repository interface {
	FindUserConnection(UserID string) (*domain.Connection, error)
	FindActiveConnection(UserID string, DomainID string) (*domain.Connection, error)
	InsertConnection(connection *domain.Connection) error
	DeleteConnection(UserID string, DeviceID string) error
	HeartbeatConnection(UserID string, DeviceID string) error
}

type requestParams struct {
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
}

type service struct {
	repository Repository
}

type Config struct {
	Repository Repository
}

func New(config Config) *service {
	return &service{
		repository: config.Repository,
	}
}

func (s *service) Connect(w http.ResponseWriter, r *http.Request) {
	var params requestParams
	if err := parseRequest(r, &params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	connection, _ := s.repository.FindUserConnection(params.UserID)
	// TODO: this also rises the error if connection is not found; need to fix this
	//if err != nil {
	//	http.Error(w, "Failed to find connection"+err.Error(), http.StatusInternalServerError)
	//	return
	//}

	if connection != nil {
		http.Error(w, "Connection already exists", http.StatusConflict)
		return
	}

	newConnection := domain.Connection{
		UserID:        params.UserID,
		DeviceID:      params.DeviceID,
		ConnectedAt:   time.Now(),
		LastHeartbeat: time.Now(),
	}

	err := s.repository.InsertConnection(&newConnection)
	if err != nil {
		http.Error(w, "Failed to insert connection"+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Connection approved; User ID: %s, Device ID: %s", params.UserID, params.DeviceID)
}

func (s *service) Disconnect(w http.ResponseWriter, r *http.Request) {
	var params requestParams
	if err := parseRequest(r, &params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := s.repository.DeleteConnection(params.UserID, params.DeviceID)
	if err != nil {
		http.Error(w, "Failed to delete connection"+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Connection deleted; User ID: %s, Device ID: %s", params.UserID, params.DeviceID)
}

func (s *service) Heartbeat(w http.ResponseWriter, r *http.Request) {
	var params requestParams
	if err := parseRequest(r, &params); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := s.repository.HeartbeatConnection(params.UserID, params.DeviceID)
	if err != nil {
		http.Error(w, "Failed to update connection"+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Heartbeat received; User ID: %s, Device ID: %s", params.UserID, params.DeviceID)
}

func parseRequest(r *http.Request, params *requestParams) error {
	return json.NewDecoder(r.Body).Decode(params)
}
