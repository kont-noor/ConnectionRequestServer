package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestConnect(t *testing.T) {
	t.Run("Successful connect", func(t *testing.T) {
		t.Parallel()

		server := newMockServer(http.StatusOK, http.StatusOK)
		defer server.Server.Close()

		client := newClient(server.Server.URL)

		err := client.Connect()
		assert.NoError(t, err)

		time.Sleep(2 * time.Second)

		server.mu.Lock()
		assert.GreaterOrEqual(t, server.HeartbeatCount, 1)
		server.mu.Unlock()
	})

	t.Run("Connection exists", func(t *testing.T) {
		t.Parallel()

		server := newMockServer(http.StatusConflict, http.StatusOK)
		defer server.Server.Close()

		client := newClient(server.Server.URL)

		err := client.Connect()
		assert.Error(t, err)

		time.Sleep(2 * time.Second)

		server.mu.Lock()
		assert.Zero(t, server.HeartbeatCount)
		server.mu.Unlock()
	})

	t.Run("Internal server error", func(t *testing.T) {
		t.Parallel()

		server := newMockServer(http.StatusInternalServerError, http.StatusOK)
		defer server.Server.Close()

		client := newClient(server.Server.URL)

		err := client.Connect()
		assert.Error(t, err)

		time.Sleep(2 * time.Second)

		server.mu.Lock()
		assert.Zero(t, server.HeartbeatCount)
		server.mu.Unlock()
	})
}

func TestDisconnect(t *testing.T) {
	t.Run("Successful disconnect", func(t *testing.T) {
		t.Parallel()

		server := newMockServer(http.StatusOK, http.StatusOK)
		defer server.Server.Close()

		client := newClient(server.Server.URL)

		client.Connect()
		err := client.Disconnect()
		assert.NoError(t, err)

		time.Sleep(2 * time.Second)

		server.mu.Lock()
		assert.Zero(t, server.HeartbeatCount)
		server.mu.Unlock()
	})

	t.Run("Failed disconnect", func(t *testing.T) {
		t.Parallel()

		server := newMockServer(http.StatusOK, http.StatusInternalServerError)
		defer server.Server.Close()

		client := newClient(server.Server.URL)

		client.Connect()
		err := client.Disconnect()
		assert.Error(t, err)

		time.Sleep(2 * time.Second)

		server.mu.Lock()
		assert.GreaterOrEqual(t, server.HeartbeatCount, 1)
		server.mu.Unlock()
	})
}

func newClient(url string) *Client {
	return New(Config{
		Host:     url,
		UserID:   "u1",
		DeviceID: "d1",
		Log:      zap.NewNop(),
	})
}
