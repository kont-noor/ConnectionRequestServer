package clientapp

import (
	"connection_request_server/internal/client"
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

	client1 := client.New(client.Config{Host: fmt.Sprintf("http://%s:%s/api/v1", serverHostname, serverPort), UserID: "15", DeviceID: "1", Log: logger})
	client2 := client.New(client.Config{Host: fmt.Sprintf("http://%s:%s/api/v1", serverHostname, serverPort), UserID: "15", DeviceID: "2", Log: logger})
	client1.Connect()
	client2.Connect()

	for {
		time.Sleep(10 * time.Second)
		client1.Disconnect()
		client2.Connect()
		client1.Connect()
		time.Sleep(10 * time.Second)
		client2.Disconnect()
		client1.Connect()
		client2.Connect()
	}
}
