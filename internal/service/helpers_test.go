package service

import (
	"bytes"
	"connection_request_server/internal/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func createRequest(method, path string, params *requestParams) *http.Request {
	var buf bytes.Buffer

	if params != nil {
		if err := json.NewEncoder(&buf).Encode(params); err != nil {
			panic("Failed to encode params: " + err.Error())
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}
func createMockService() (*service, *mockRepository) {
	repo := &mockRepository{
		Connections: []*domain.Connection{
			{
				ID:       "CONNECTION_ID",
				UserID:   "USER_ID",
				DeviceID: "DEVICE_ID",
			},
		},
	}
	service := New(Config{Repository: repo})

	return service, repo
}
