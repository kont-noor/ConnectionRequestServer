package service

import (
	"bytes"
	"connection_request_server/internal/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	connections := []*domain.Connection{
		{
			ID:       "CONNECTION_ID",
			UserID:   "USER_ID",
			DeviceID: "DEVICE_ID",
		},
	}
	repo := &MockRepository{
		Connections: connections,
	}
	service := New(Config{Repository: repo})

	t.Run("Bad Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/connect", nil)
		res := httptest.NewRecorder()

		service.Connect(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, len(repo.Connections), 1)
	})

	t.Run("Connection exists", func(t *testing.T) {
		params := requestParams{UserID: "USER_ID"}
		jsonBody, err := json.Marshal(params)

		if err != nil {
			t.Fatalf("Failed to marshal request params: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/connect", bytes.NewReader(jsonBody))
		res := httptest.NewRecorder()

		service.Connect(res, req)

		assert.Equal(t, http.StatusConflict, res.Code)
		assert.Equal(t, len(repo.Connections), 1)
	})

	t.Run("Create connection", func(t *testing.T) {
		params := requestParams{UserID: "User1", DeviceID: "Device1"}
		jsonBody, err := json.Marshal(params)

		if err != nil {
			t.Fatalf("Failed to marshal request params: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/connect", bytes.NewReader(jsonBody))
		res := httptest.NewRecorder()

		service.Connect(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, len(repo.Connections), 2)

		connection := repo.Connections[len(repo.Connections)-1]

		assert.Equal(t, "User1", connection.UserID)
		assert.Equal(t, "Device1", connection.DeviceID)
		assert.NotNil(t, connection.ConnectedAt)
		assert.NotNil(t, connection.LastHeartbeat)
	})
}
