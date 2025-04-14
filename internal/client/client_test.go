package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestConnect(t *testing.T) {
	tests := []struct {
		Name             string
		ReceivedResponse int
		IsError          bool
		Heartbeat        bool
	}{
		{
			Name:             "Server returns 500",
			ReceivedResponse: http.StatusInternalServerError,
			IsError:          true,
			Heartbeat:        false,
		},
		{
			Name:             "Connection already exists",
			ReceivedResponse: http.StatusConflict,
			IsError:          true,
			Heartbeat:        false,
		},
		{
			Name:             "Connected successfully",
			ReceivedResponse: http.StatusOK,
			IsError:          false,
			Heartbeat:        true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			log := zap.NewNop()

			server := newMockServer(test.ReceivedResponse, http.StatusOK)
			defer server.Server.Close()

			client := New(Config{
				Host:     server.Server.URL,
				UserID:   "u1",
				DeviceID: "d1",
				Log:      log,
			})

			err := client.Connect()
			if test.IsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			time.Sleep(2 * time.Second)

			server.mu.Lock()
			if test.Heartbeat {
				assert.GreaterOrEqual(t, server.HeartbeatCount, 1)
			} else {
				assert.Zero(t, server.HeartbeatCount)
			}
			server.mu.Unlock()
		})
	}
}

func TestDisconnect(t *testing.T) {
	tests := []struct {
		Name             string
		ReceivedResponse int
		IsError          bool
		Heartbeat        bool
	}{
		{
			Name:             "Success",
			ReceivedResponse: http.StatusOK,
			IsError:          false,
			Heartbeat:        false,
		},
		{
			Name:             "Fail",
			ReceivedResponse: http.StatusInternalServerError,
			IsError:          true,
			Heartbeat:        true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			log := zap.NewNop()

			server := newMockServer(http.StatusOK, test.ReceivedResponse)
			defer server.Server.Close()

			client := New(Config{
				Host:     server.Server.URL,
				UserID:   "u1",
				DeviceID: "d1",
				Log:      log,
			})

			client.Connect()
			err := client.Disconnect()
			if test.IsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			time.Sleep(2 * time.Second)

			server.mu.Lock()
			if test.Heartbeat {
				assert.GreaterOrEqual(t, server.HeartbeatCount, 1)
			} else {
				assert.Zero(t, server.HeartbeatCount)
			}
			server.mu.Unlock()
		})
	}
}
