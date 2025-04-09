package clientapp

import (
	"connection_request_server/internal/client"
	randomticker "connection_request_server/pkg/random_ticker"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

func Run() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Init client app")
	time.Sleep(4 * time.Second)

	serverHostname := os.Getenv("SERVER_HOSTNAME")
	serverPort := os.Getenv("SERVER_PORT")
	userId := os.Getenv("USER_ID")
	deviceId := os.Getenv("DEVICE_ID")

	connected := false
	client1 := client.New(client.Config{Host: fmt.Sprintf("http://%s:%s/api/v1", serverHostname, serverPort), UserID: userId, DeviceID: deviceId, Log: logger})
	err := client1.Connect()

	if err == nil {
		connected = true
	}

	ticker := randomticker.New(1*time.Second, 9*time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if connected {
			err = client1.Disconnect()
			if err == nil {
				connected = false
			}
		} else {
			err = client1.Connect()
			if err == nil {
				connected = true
			}
		}
	}
}
