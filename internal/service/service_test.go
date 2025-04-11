package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tests := []struct {
		name           string
		input          *requestParams
		expectedStatus int
		expectedConns  int
		validate       func(*testing.T, *mockRepository)
	}{
		{
			name:           "Bad Request - no body",
			input:          nil,
			expectedStatus: http.StatusBadRequest,
			expectedConns:  1,
		},
		{
			name:           "Connection exists",
			input:          &requestParams{UserID: "USER_ID", DeviceID: "DEVICE_ID"},
			expectedStatus: http.StatusConflict,
			expectedConns:  1,
		},
		{
			name:           "Create connection",
			input:          &requestParams{UserID: "User1", DeviceID: "Device1"},
			expectedStatus: http.StatusOK,
			expectedConns:  2,
			validate: func(t *testing.T, repo *mockRepository) {
				conn := repo.Connections[len(repo.Connections)-1]
				assert.Equal(t, "User1", conn.UserID)
				assert.Equal(t, "Device1", conn.DeviceID)
				assert.NotZero(t, conn.ConnectedAt)
				assert.NotZero(t, conn.LastHeartbeat)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, repo := createMockService()
			req := createRequest(http.MethodPost, "/connect", test.input)
			res := httptest.NewRecorder()

			service.Connect(res, req)

			assert.Equal(t, test.expectedStatus, res.Code)
			assert.Len(t, repo.Connections, test.expectedConns)

			if test.validate != nil {
				test.validate(t, repo)
			}
		})
	}
}

func TestDisconnect(t *testing.T) {
	tests := []struct {
		name           string
		input          *requestParams
		expectedStatus int
		expectedConns  int
		validate       func(*testing.T, *mockRepository)
	}{
		{
			name:           "Bad Request - no body",
			input:          nil,
			expectedStatus: http.StatusBadRequest,
			expectedConns:  1,
		},
		{
			name:           "Connection does not exist",
			input:          &requestParams{UserID: "User1", DeviceID: "Device1"},
			expectedStatus: http.StatusInternalServerError,
			expectedConns:  1,
		},
		{
			name:           "Delete connection",
			input:          &requestParams{UserID: "USER_ID", DeviceID: "DEVICE_ID"},
			expectedStatus: http.StatusOK,
			expectedConns:  0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, repo := createMockService()
			req := createRequest(http.MethodDelete, "/disconnect", test.input)
			res := httptest.NewRecorder()

			service.Disconnect(res, req)

			assert.Equal(t, test.expectedStatus, res.Code)
			assert.Len(t, repo.Connections, test.expectedConns)

			if test.validate != nil {
				test.validate(t, repo)
			}
		})
	}
}

func TestHeartbeat(t *testing.T) {
	tests := []struct {
		name             string
		input            *requestParams
		expectedStatus   int
		heartbeatUpdated bool
	}{
		{
			name:             "Bad Request - no body",
			input:            nil,
			expectedStatus:   http.StatusBadRequest,
			heartbeatUpdated: false,
		},
		{
			name:             "Connection does not exist",
			input:            &requestParams{UserID: "User1", DeviceID: "Device1"},
			expectedStatus:   http.StatusInternalServerError,
			heartbeatUpdated: false,
		},
		{
			name:             "Update heartbeat",
			input:            &requestParams{UserID: "USER_ID", DeviceID: "DEVICE_ID"},
			expectedStatus:   http.StatusOK,
			heartbeatUpdated: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, repo := createMockService()
			req := createRequest(http.MethodPut, "/heartbeat", test.input)
			res := httptest.NewRecorder()

			service.Heartbeat(res, req)

			assert.Equal(t, test.expectedStatus, res.Code)
			if test.heartbeatUpdated {
				assert.NotEqual(t, time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), repo.Connections[0].LastHeartbeat)
			} else {
				assert.Equal(t, time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC), repo.Connections[0].LastHeartbeat)
			}
		})
	}
}
