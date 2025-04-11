package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tests := []struct {
		name           string
		input          *requestParams
		expectedStatus int
		expectedConns  int
		validate       func(*testing.T, *MockRepository)
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
			validate: func(t *testing.T, repo *MockRepository) {
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
			req := createRequest(http.MethodGet, "/connect", test.input)
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
