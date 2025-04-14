package router

import (
	"bytes"
	"net/http"
	"testing"

	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRouter(t *testing.T) {
	mock := newMockHandler()
	log := zap.NewNop()

	r := New(Config{
		APIHandlers: mock,
		Log:         log,
	})

	tests := []struct {
		path     string
		expected string
		method   string
	}{
		{"/api/v1/connect", "connected", http.MethodPost},
		{"/api/v1/disconnect", "disconnected", http.MethodDelete},
		{"/api/v1/heartbeat", "heartbeat", http.MethodPut},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(`{"user_id":"u","device_id":"d"}`))
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, tt.expected, rec.Body.String())
		assert.True(t, mock.called[tt.expected], "Handler not called for "+tt.path)
	}
}
