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
		ReceivedResponce int
		IsError          bool
		Heartbeat        bool
	}{
		{
			Name:             "Server returns 500",
			ReceivedResponce: http.StatusInternalServerError,
			IsError:          true,
			Heartbeat:        false,
		},
		{
			Name:             "Connection already exists",
			ReceivedResponce: http.StatusConflict,
			IsError:          true,
			Heartbeat:        false,
		},
		{
			Name:             "Connected successfully",
			ReceivedResponce: http.StatusOK,
			IsError:          false,
			Heartbeat:        true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			log := zap.NewNop()

			server := newMockServer(test.ReceivedResponce, http.StatusOK)
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

			if test.Heartbeat {
				assert.GreaterOrEqual(t, server.HeartbeatCount, 1)
			} else {
				assert.Zero(t, server.HeartbeatCount)
			}
		})
	}
}
